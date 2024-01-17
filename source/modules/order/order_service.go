package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	DBConnection "sorteio-api/source/database"
	"sorteio-api/source/dto"
	"sorteio-api/source/models"
	"sorteio-api/source/modules/sorteio"
	"sorteio-api/source/modules/ticket"
	"sorteio-api/source/modules/user"
	efi "sorteio-api/source/resources/efi_sdk"
	"sorteio-api/source/utils"
	"strconv"
	"time"

	"github.com/gerencianet/gn-api-sdk-go/gerencianet/pix"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SERVICE
type IOrderService interface {
	GetAllOrders(page int, pageSize int, filter bson.M) ([]models.Order, int, error)
	GetAllOrdersFromSorteioId(sorteioId int64, page int, pageSize int) ([]models.Order, int, error)
	GetAOrderFromTransactionId(transacionId string) (models.Order, bool)
	GetAOrderFromObjId(orderId string) (models.Order, bool)
	GetAOrderFromId(orderId string) (models.Order, bool)
	CreatorOrderSteps(data dto.OrderRequestDto) (bool, models.Order)
	OrderConfirmPayment(TransactionId string) (dto.OrderResponseDto, bool)
	CreateTicketsFromOrderConfirm(sorteio models.Sorteio, orderId primitive.ObjectID, userId primitive.ObjectID, transactionId string, ticketsAmount int) []int
	UpdateOrder(data models.Order)
	GetOrdersFromUserCredential(credential string, value string) ([]models.Order, int, error)
	GetOrdersFromUserId(userId primitive.ObjectID) ([]models.Order, int, error)
	CreateCharge(dueSeconds int, cpf string, name string, price string, orderId string) string
	GenerateQRCode(transactionId string) string
	GetUsernameFromUserId(userId primitive.ObjectID) string
	AnalyzePendingOrders()
}

type OrderService struct {
	OrderRepository IOrderRepository
	Validate        *validator.Validate
	UserService     user.IUserService
	SorteioService  sorteio.ISorteioService
	TicketService   ticket.ITicketService
}

func InitOrderService(OrderRepository IOrderRepository, validate *validator.Validate, UserService user.IUserService, SorteioService sorteio.ISorteioService, TicketService ticket.ITicketService) IOrderService {
	return &OrderService{
		OrderRepository: OrderRepository,
		Validate:        validate,
		UserService:     UserService,
		SorteioService:  SorteioService,
		TicketService:   TicketService,
	}
}

func (s *OrderService) GetAllOrders(page int, pageSize int, filter bson.M) ([]models.Order, int, error) {
	results, err, ctx := DBConnection.FindMultipleDocuments("order", filter, page, pageSize)

	count := results.RemainingBatchLength()

	var Orders []models.Order

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleOrder models.Order
		if err = results.Decode(&singleOrder); err != nil {
			panic(err)
		}

		Orders = append(Orders, singleOrder)
	}

	return Orders, count, err
}

func (s *OrderService) GetAllOrdersFromSorteioId(sorteioId int64, page int, pageSize int) ([]models.Order, int, error) {
	filter := bson.M{"sorteio_id": sorteioId}
	results, err, ctx := DBConnection.FindMultipleDocuments("order", filter, page, pageSize)

	count := results.RemainingBatchLength()

	var Orders []models.Order

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleOrder models.Order
		if err = results.Decode(&singleOrder); err != nil {
			panic(err)
		}

		Orders = append(Orders, singleOrder)
	}

	return Orders, count, err
}

func (s *OrderService) GetAOrderFromTransactionId(transacionId string) (models.Order, bool) {
	filter := bson.M{"transaction_id": transacionId}
	result, _ := DBConnection.FindADocument("order", filter)

	var order models.Order
	err := result.Decode(&order)

	if err != nil {
		return models.Order{}, false
	}

	return order, true
}

func (s *OrderService) GetAOrderFromObjId(orderId string) (models.Order, bool) {
	var order models.Order

	objId, _ := primitive.ObjectIDFromHex(orderId)

	result, _ := DBConnection.FindADocument("order", bson.M{"_id": objId})
	err := result.Decode(&order)

	if err != nil {
		return models.Order{}, false
	}

	return order, true
}

func (s *OrderService) GetAOrderFromId(orderId string) (models.Order, bool) {
	ordId, _ := strconv.Atoi(orderId)

	filter := bson.M{"order_id": ordId}
	result, _ := DBConnection.FindADocument("order", filter)

	var order models.Order
	err := result.Decode(&order)

	if err != nil {
		return models.Order{}, false
	}

	return order, true
}

// FLUXO PARA CRIAR UMA ORDER:

// 1 Criar um usuário
// 2 Gerar os tickets
// 3 Criar uma order
// 4 Fazer o Request do Pagamento
// 5 Aguardar o Status do pagamento
// 6 Retornar ao usuário os tickets com o sucesso do pagamento.

func (s *OrderService) CreatorOrderSteps(data dto.OrderRequestDto) (bool, models.Order) {
	orderId := primitive.NewObjectID()

	//// FIRST STEP -- USER
	var Credentials = models.Credentials{
		Email: data.Email,
		CPF:   data.Cpf,
		Phone: data.Phone,
	}

	user := s.UserService.GetUserFromCredentials(Credentials)

	if user.ID.IsZero() {

		userRoles := []models.Role{models.Role{
			Name: "user",
			Role: 0,
		}}

		newUserID := primitive.NewObjectID()

		newUser := models.User{
			ID:      newUserID,
			Name:    data.Name,
			Email:   data.Email,
			Phone:   data.Phone,
			CPF:     data.Cpf,
			Blocked: false,
			Roles:   userRoles,
		}

		s.UserService.CreateUser(newUser)
		user.ID = newUserID
	}

	//// SEARCH SORTEIO FROM ID
	sortId := strconv.Itoa(data.Sorteio)
	sorteio, sorteioErr := s.SorteioService.GetASorteioFromId(sortId)

	if sorteio.TicketsSold+data.Tickets > sorteio.Tickets.Amount {
		panic(errors.New("Limite máximo de bilhetes atingido"))
	}

	if !sorteioErr {
		panic(sorteioErr)
	}

	//// SECOND STEP -- TICKETS
	ordersCount, _ := DBConnection.CollectionCount("order")

	var secondsToExpire int = (sorteio.TimeExpireOrder * 60)
	var expireAt = time.Now().Add(time.Second * time.Duration(secondsToExpire))

	//// CALCULE TOTAL PRICE
	var totalValue = (sorteio.Tickets.Price * data.Tickets)

	s.SorteioService.UpdateTicketsSold(strconv.Itoa(int(sorteio.SorteioId)), data.Tickets)

	// if totalValue >= sorteio.Discount.Minimal {
	// 	totalValue = totalValue - (sorteio.Discount.Amount * data.Tickets)
	// }

	orderNumber := ordersCount + 1

	stringPrice := utils.PriceToFloatString(uint32(totalValue))

	// fmt.Println("ESSE é o preço ::", stringPrice)

	chargeRes := s.CreateCharge(secondsToExpire, data.Cpf, data.Name, stringPrice, fmt.Sprintf("%s", orderNumber))

	var charge map[string]interface{}

	if err := json.Unmarshal([]byte(chargeRes), &charge); err != nil {
		panic(err)
	}

	loc := charge["loc"].(map[string]interface{})

	var txid string = fmt.Sprintf("%s", charge["txid"])

	qrCodeRes := s.GenerateQRCode(fmt.Sprintf("%v", loc["id"]))
	var qrCode map[string]interface{}

	if errQrCode := json.Unmarshal([]byte(qrCodeRes), &qrCode); errQrCode != nil {
		panic(errQrCode)
	}

	var QRCodeImage string = fmt.Sprintf("%s", qrCode["imagemQrcode"])
	var QRCodeString string = fmt.Sprintf("%s", qrCode["qrcode"])

	var transactionId string = txid

	//// THIRD STEP -- CREATE ORDER
	newOrder := models.Order{
		ID:            orderId,
		OrderId:       orderNumber,
		UserId:        user.ID,
		UserName:      user.Name,
		SorteioId:     sorteio.SorteioId,
		TransactionId: transactionId,

		PaymentMethod: "pix",
		// Tickets:       tickets,
		TicketsAmount: data.Tickets,
		QRCodeImage:   QRCodeImage,
		QRCodeString:  QRCodeString,

		Status:   "pending",
		Total:    totalValue,
		ExpireAt: expireAt,
		Referal:  data.Referal,

		CreatedAt: time.Now(),
	}

	resultId, status := s.CreateOrder(newOrder)

	if status {
		utils.DebugPrint(" resultId :: ", resultId)
	}

	// orderDto := dto.OrderResponseDto{
	// 	OrderId:   orderNumber,
	// 	SorteioId: sorteio.ID.String(),
	// 	// Tickets:       tickets,
	// 	TicketsAmount: data.Tickets,
	// 	TransaciontId: transactionId,
	// 	Total:         totalValue,
	// 	ExpireAt:      expireAt.Unix(),
	// }

	return status, newOrder
}

func (s *OrderService) OrderConfirmPayment(TransactionId string) (dto.OrderResponseDto, bool) {
	order, status := s.GetAOrderFromTransactionId(TransactionId)

	sortId := strconv.Itoa(int(order.SorteioId))
	sorteio, _ := s.SorteioService.GetASorteioFromId(sortId)

	tickets := s.CreateTicketsFromOrderConfirm(sorteio, order.ID, order.UserId, order.TransactionId, order.TicketsAmount)

	order.Tickets = tickets
	order.Status = "paid"

	orderDto := dto.OrderResponseDto{
		OrderId:       order.OrderId,
		SorteioId:     sorteio.ID.String(),
		Tickets:       tickets,
		TicketsAmount: order.TicketsAmount,
		TransaciontId: order.TransactionId,
		Status:        "paid",
		Total:         order.Total,
		ExpireAt:      order.ExpireAt.Unix(),
	}

	s.UpdateOrder(order)

	return orderDto, status
}

func (s *OrderService) CreateTicketsFromOrderConfirm(sorteio models.Sorteio, orderId primitive.ObjectID, userId primitive.ObjectID, transactionId string, ticketsAmount int) []int {
	var tickets = []int{}

	ticketsToCreate := []interface{}{}

	var sum int = 0

	for sum < ticketsAmount {
		ticketNumber := s.TicketService.GenerateTicketNumber(sorteio.SorteioId, sorteio.Tickets.Amount)

		if ticketNumber != -1 {
			newTicket := models.Ticket{
				ID:            primitive.NewObjectID(),
				Number:        ticketNumber,
				UserId:        userId,
				SorteioId:     sorteio.SorteioId,
				OrderId:       orderId,
				Status:        true,
				CreatedAt:     time.Now(),
				TransactionId: transactionId,
			}

			ticketsToCreate = append(ticketsToCreate, newTicket)
			tickets = append(tickets, ticketNumber)
			s.TicketService.CreateTicket(newTicket)

			sum += 1
		}
	}

	// s.TicketService.CreateManyTickets(ticketsToCreate)

	return tickets
}

func (s *OrderService) UpdateOrder(data models.Order) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("order")
	defer cancel()

	if validationErr := s.Validate.Struct(&data); validationErr != nil {
		panic(validationErr)
	}

	filter := bson.M{"order_id": data.OrderId}

	update := bson.M{
		"$set": bson.M{
			"tickets":  data.Tickets,
			"status":   data.Status,
			"rejected": data.Rejected,
		},
	}

	_, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		panic(err)
	}
}

func (s *OrderService) CreateOrder(data models.Order) (interface{}, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("order")
	defer cancel()

	if validationErr := s.Validate.Struct(&data); validationErr != nil {
		panic(validationErr)
	}

	result, err := coll.InsertOne(ctx, data)

	if err != nil {
		panic(err)
	}

	if result != nil {
		return nil, false
	}

	return result.InsertedID, true
}

func (s *OrderService) GetOrdersFromUserCredential(credential string, value string) ([]models.Order, int, error) {
	user, err := s.UserService.GetUserFromCredential(credential, value)

	Orders, count, _ := s.GetOrdersFromUserId(user.ID)

	return Orders, count, err
}

func (s *OrderService) GetUsernameFromUserId(userId primitive.ObjectID) string {
	user, _ := s.UserService.GetUserFromId(userId)

	//	if !result {
	//		panic(errors.New("Usuário não encontrado"));
	//	}

	return user.Name
}

func (s *OrderService) GetOrdersFromUserId(userId primitive.ObjectID) ([]models.Order, int, error) {
	filter := bson.M{"user_id": userId}

	results, err, ctx := DBConnection.FindMultipleDocuments("order", filter, 0, 0)
	count := results.RemainingBatchLength()

	var Orders []models.Order

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleOrder models.Order
		if err = results.Decode(&singleOrder); err != nil {
			panic(err)
		}

		Orders = append(Orders, singleOrder)
	}

	return Orders, count, err
}

func (s *OrderService) OrdersNotPaid() []models.Order {
	filter := bson.M{"status": "pending"}

	response, _, _ := s.GetAllOrders(0, 0, filter)

	var Orders []models.Order

	for _, order := range response {
		if time.Now().Unix() >= order.ExpireAt.Unix() {
			order.Status = "expired"
			s.UpdateOrder(order)
			s.SorteioService.UpdateTicketsSold(strconv.Itoa(int(order.SorteioId)), -(order.TicketsAmount))
		}
	}

	return Orders
}

func (s *OrderService) AnalyzePendingOrders() {
	// tickTime := time.Now().Unix()

	s.OrdersNotPaid()

	time.Sleep(60 * time.Second)

	s.AnalyzePendingOrders()
}

func (s *OrderService) CreateCharge(dueSeconds int, cpf string, name string, price string, orderId string) string {
	GN := pix.NewGerencianet(efi.Credentials)

	chavePix := os.Getenv("EFI_CHAVE_PIX")
	efiPaymentFull := os.Getenv("EFI_PAYMENT_FULL")

	if efiPaymentFull == "false" {
		price = "0.01"
	}

	utils.DebugPrint(efi.Credentials)

	body := map[string]interface{}{

		"calendario": map[string]interface{}{
			"expiracao": dueSeconds,
		},
		"devedor": map[string]interface{}{
			"cpf":  utils.NormalizeCPF(cpf),
			"nome": name,
		},
		"valor": map[string]interface{}{
			"original": price,
		},
		"chave":              chavePix,
		"solicitacaoPagador": "Sorteio - Cotas",
	}

	res, err := GN.CreateImmediateCharge(body)

	if err != nil {
		panic(err)
	}

	return res
}

func (s *OrderService) GenerateQRCode(transactionId string) string {
	GN := pix.NewGerencianet(efi.Credentials)
	res, err := GN.PixGenerateQRCode(transactionId)

	if err != nil {
		panic(err)
	}

	return res
}
