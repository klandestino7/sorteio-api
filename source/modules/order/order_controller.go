package order

import (
	"net/http"
	"sorteio-api/source/dto"
	gnEvent "sorteio-api/source/modules/gn_event"
	efi "sorteio-api/source/resources/efi_sdk"
	"strconv"

	"os"

	"github.com/gerencianet/gn-api-sdk-go/gerencianet/pix"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// CONTROLLER
type IOrderController interface {
	RequestOrderStatusFromTransactionId(c *gin.Context)
	RequestOrderFromObjId(c *gin.Context)
	RequestOrderFromId(c *gin.Context)
	RequestAllOrders(c *gin.Context)
	RequestOrdersPaid(c *gin.Context)
	RequestOrdersFromSorteioId(c *gin.Context)
	RequestOrdersWithUsers(c *gin.Context)
	RequestOrdersFromCredential(c *gin.Context)
	RequestCreateOrder(c *gin.Context)
	PixHandle(c *gin.Context)
	ConfigWebhookRequest(c *gin.Context)
}

type OrderController struct {
	OrderService   IOrderService
	GnEventService gnEvent.IGnEventService
}

func InitOrderController(OrderService IOrderService, GnEvent gnEvent.IGnEventService) IOrderController {
	return &OrderController{
		OrderService:   OrderService,
		GnEventService: GnEvent,
	}
}

func (ct *OrderController) RequestOrderStatusFromTransactionId(c *gin.Context) {
	transactionId := c.Param("transactionId")

	if transactionId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "transactionId Invalid"})
		return
	}

	response, status := ct.OrderService.GetAOrderFromTransactionId(transactionId)

	if !status {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "transactionId Invalid"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": response.Status})
}

func (ct *OrderController) RequestOrderFromObjId(c *gin.Context) {
	orderId := c.Param("orderId")

	if orderId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "orderId Invalid"})
		return
	}

	response, status := ct.OrderService.GetAOrderFromObjId(orderId)

	if !status {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "order Invalid"})
		return
	}

	order := dto.CreateOrderResponse(&response)

	c.IndentedJSON(http.StatusOK, gin.H{"order": order})
}

func (ct *OrderController) RequestOrderFromId(c *gin.Context) {
	orderId := c.Param("orderId")

	if orderId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "orderId Invalid"})
		return
	}

	response, status := ct.OrderService.GetAOrderFromId(orderId)

	if !status {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "order Invalid"})
		return
	}

	order := dto.CreateOrderResponse(&response)

	c.IndentedJSON(http.StatusOK, gin.H{"order": order})
}

func (ct *OrderController) RequestAllOrders(c *gin.Context) {
	filter := bson.M{}

	response, count, err := ct.OrderService.GetAllOrders(0, 0, filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		return
	}

	var Orders []map[string]interface{}
	for _, order := range response {
		Orders = append(Orders, dto.CreateOrderResponse(&order))
	}

	c.IndentedJSON(http.StatusOK, gin.H{"orders": Orders, "count": count})
}

func (ct *OrderController) RequestOrdersPaid(c *gin.Context) {
	filter := bson.M{"status": 1}

	response, count, err := ct.OrderService.GetAllOrders(0, 0, filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		return
	}

	var Orders []map[string]interface{}
	for _, order := range response {
		Orders = append(Orders, dto.CreateOrderResponse(&order))
	}

	c.IndentedJSON(http.StatusOK, gin.H{"orders": Orders, "count": count})
}

func (ct *OrderController) RequestOrdersFromSorteioId(c *gin.Context) {
	sorteioId := c.Query("sorteioId")
	page := c.DefaultQuery("page", "0")

	sorteioIdInt64, err := strconv.ParseInt(sorteioId, 10, 64)
	pageInt, err := strconv.Atoi(page)

	response, count, err := ct.OrderService.GetAllOrdersFromSorteioId(sorteioIdInt64, pageInt, 50)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		return
	}

	var Orders []map[string]interface{}
	for _, order := range response {
		order.UserName = ct.OrderService.GetUsernameFromUserId(order.UserId)
		Orders = append(Orders, dto.CreateOrderResponse(&order))
	}

	c.IndentedJSON(http.StatusOK, gin.H{"orders": Orders, "count": count})
}

func (ct *OrderController) RequestOrdersWithUsers(c *gin.Context) {
	filter := bson.M{}
	response, count, err := ct.OrderService.GetAllOrders(0, 0, filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		return
	}

	var Orders []map[string]interface{}
	for _, order := range response {
		Orders = append(Orders, dto.CreateOrderResponse(&order))
	}

	c.IndentedJSON(http.StatusOK, gin.H{"orders": Orders, "count": count})
}

func (ct *OrderController) RequestOrdersFromCredential(c *gin.Context) {
	var json dto.OrderRequestFromCredential

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		panic(err)
	}

	var Orders []map[string]interface{}

	result, count, err := ct.OrderService.GetOrdersFromUserCredential(json.Credential, json.Value)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"orders": Orders, "count": count})
	}

	for _, order := range result {
		Orders = append(Orders, dto.CreateOrderResponse(&order))
	}

	c.JSON(http.StatusOK, gin.H{"orders": Orders, "count": count})
}

func (ct *OrderController) RequestCreateOrder(c *gin.Context) {
	var json dto.OrderRequestDto

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		panic(err)
	}

	_, order := ct.OrderService.CreatorOrderSteps(json)

	// if result {

	orderResponse := dto.CreateOrderResponse(&order)

	c.JSON(http.StatusOK, gin.H{"order": orderResponse})
	// }
}

func (ct *OrderController) ConfigWebhookRequest(c *gin.Context) {
	credentials := efi.Credentials

	// fmt.Println("Credentials :: ", credentials)

	GN := pix.NewGerencianet(credentials)

	chavePix := os.Getenv("EFI_CHAVE_PIX")
	efiWebhookReturn := os.Getenv("EFI_WEBHOOK_RETURN")

	body := map[string]interface{}{
		"webhookUrl": efiWebhookReturn,
	}

	res, _ := GN.PixConfigWebhook(chavePix, body)

	c.JSON(http.StatusOK, res)
}

func (ct *OrderController) PixHandle(c *gin.Context) {
	var json dto.PixHandleDto

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		panic(err)
	}

	for _, pix := range json.Pix {
		ct.OrderService.OrderConfirmPayment(pix.Txid)
		ct.GnEventService.CreateGnEvent(pix)
	}

	c.JSON(http.StatusOK, "done")
}

// func ConfirmPix(c *gin.Context) {
// 	service.OrderConfirmPayment("c1e0ef53fe8541d4a782afa17751afd6")

// 	c.JSON(http.StatusOK, "Olá mundo")
// }

func (ct *OrderController) WebhookRequest(c *gin.Context) {
	c.JSON(http.StatusOK, "Olá mundo")
}
