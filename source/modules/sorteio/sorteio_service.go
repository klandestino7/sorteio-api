package sorteio

import (
	"context"
	"errors"
	"fmt"
	DBConnection "sorteio-api/source/database"
	"sorteio-api/source/dto"
	"sorteio-api/source/models"
	"sorteio-api/source/modules/ticket"
	"sorteio-api/source/modules/winner"
	"sorteio-api/source/utils"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SERVICE
type ISorteioService interface {
	GetAllSorteios(page int, pageSize int) ([]models.Sorteio, int, error)
	GetFullSorteios(page int, pageSize int) ([]models.Sorteio, int, error)
	SorteioEarning(sorteioId primitive.ObjectID) int
	SorteioTicketsPurchased(sorteioId int64) int
	CreateSorteio(data dto.SorteioCreateRequestDto) (interface{}, error)
	GetASorteioFromId(sorteioId string) (models.Sorteio, bool)
	GetASorteioFromObjId(sorteioObjId string) (models.Sorteio, bool)
	GenerateAWinnerToSorteio(sorteioObjId string) models.Ticket
	UpdateSorteio(data dto.UpdateSorteioDto)
	UpdateTicketsSold(sorteioId string, sold int)
}

type SorteioService struct {
	SorteioRepository ISorteioRepository
	Validate          *validator.Validate
	TicketService     ticket.ITicketService
	WinnerService     winner.IWinnerService
}

func InitSorteioService(SorteioRepository ISorteioRepository, validate *validator.Validate, TicketService ticket.ITicketService, WinnerService winner.IWinnerService) ISorteioService {
	return &SorteioService{
		SorteioRepository: SorteioRepository,
		Validate:          validate,
		TicketService:     TicketService,
		WinnerService:     WinnerService,
	}
}

func (s *SorteioService) GetAllSorteios(page int, pageSize int) ([]models.Sorteio, int, error) {
	filter := bson.M{}

	results, err, ctx := DBConnection.FindMultipleDocuments("sorteio", filter, page, pageSize)

	count := results.RemainingBatchLength()

	var Sorteios []models.Sorteio

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleSorteio models.Sorteio
		if err = results.Decode(&singleSorteio); err != nil {
			panic(err)
		}

		Sorteios = append(Sorteios, singleSorteio)
	}

	return Sorteios, count, err
}

func (s *SorteioService) GetFullSorteios(page int, pageSize int) ([]models.Sorteio, int, error) {
	filter := bson.M{}

	results, err, ctx := DBConnection.FindMultipleDocuments("sorteio", filter, page, pageSize)

	count := results.RemainingBatchLength()

	var Sorteios []models.Sorteio

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleSorteio models.Sorteio
		if err = results.Decode(&singleSorteio); err != nil {
			panic(err)
		}

		var maxTicketsAmount = singleSorteio.Tickets.Amount
		var ticketsSold = singleSorteio.TicketsSold

		var diference = float64(ticketsSold) / float64(maxTicketsAmount)
		var calcResult = float64(diference * 100)

		singleSorteio.Percentage = calcResult

		singleSorteio.Earning = s.SorteioEarning(singleSorteio.ID)

		Sorteios = append(Sorteios, singleSorteio)
	}

	return Sorteios, count, err
}

func (s *SorteioService) SorteioEarning(sorteioId primitive.ObjectID) int {
	filter := bson.M{"sorteio_id": sorteioId}

	results, err, ctx := DBConnection.FindMultipleDocuments("order", filter, 0, 0)

	var Earning int = 0

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleOrder models.Order
		if err = results.Decode(&singleOrder); err != nil {
			panic(err)
		}

		Earning += Earning + singleOrder.Total
	}

	return Earning
}

func (s *SorteioService) SorteioTicketsPurchased(sorteioId int64) int {
	filter := bson.M{"sorteio_id": sorteioId}

	results, _, _ := DBConnection.FindMultipleDocuments("ticket", filter, 0, 0)
	count := results.RemainingBatchLength()

	fmt.Println("count :: ", count)

	return count
}

func (s *SorteioService) CreateSorteio(data dto.SorteioCreateRequestDto) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("sorteio")
	defer cancel()

	utils.DebugPrint("create sorteio :: ", data)

	if validationErr := s.Validate.Struct(&data); validationErr != nil {
		panic(validationErr)
	}

	var sorteioStatus bool

	if data.Status == "true" {
		sorteioStatus = true
	} else {
		sorteioStatus = false
	}

	var displayFinishDateStatus bool

	if data.DisplayFinishDate == "true" {
		displayFinishDateStatus = true
	} else {
		displayFinishDateStatus = false
	}

	sorteiosCount, _ := DBConnection.CollectionCount("sorteio")

	ticketPrice, _ := strconv.Atoi(data.TicketsPrice)
	ticketAmount, _ := strconv.Atoi(data.TicketsAmount)
	minimalForOrder, _ := strconv.Atoi(data.MinimalTicketForOrder)
	maximumForOrder, _ := strconv.Atoi(data.MaximumTicketForOrder)

	newTickets := models.Tickets{
		Amount:          ticketAmount,
		Price:           ticketPrice,
		MinimalForOrder: minimalForOrder,
		MaximumForOrder: maximumForOrder,
	}

	var newDiscount models.Discount

	discountAmount, _ := strconv.Atoi(data.DiscountAmount)
	ticketsMinimalToDiscount, _ := strconv.Atoi(data.TicketsMinimalToDiscount)
	// limitTimetoExpireOrder, _ := strconv.Atoi(data.LimitTimeToExpireOrder)
	finishDate, _ := strconv.ParseInt(data.FinishDate, 0, 64)

	if discountAmount > 0 {
		newDiscount = models.Discount{
			Amount:  discountAmount,
			Minimal: ticketsMinimalToDiscount,
		}
	}

	var Sorteio = models.Sorteio{
		ID:          primitive.NewObjectID(),
		Name:        data.Name,
		Description: data.Description,

		// Sorteios: [],
		Tickets:           newTickets,
		Discount:          newDiscount,
		SorteioId:         sorteiosCount + 1,
		FinishDate:        time.Unix(finishDate, 0),
		DisplayFinishDate: displayFinishDateStatus,

		TimeExpireOrder: 60,
		Status:          sorteioStatus,
		CreatedAt:       time.Now(),
	}

	result, err := coll.InsertOne(ctx, Sorteio)

	if err != nil {
		return nil, err
	}

	if result != nil {
		return nil, err
	}

	return result.InsertedID, err
}

func (s *SorteioService) GetASorteioFromId(sorteioId string) (models.Sorteio, bool) {
	sortId, _ := strconv.Atoi(sorteioId)

	filter := bson.M{"sorteio_id": sortId}
	result, _ := DBConnection.FindADocument("sorteio", filter)

	var sorteio models.Sorteio
	err := result.Decode(&sorteio)

	var maxTicketsAmount = sorteio.Tickets.Amount
	var ticketsSold = sorteio.TicketsSold

	var diference = float64(ticketsSold) / float64(maxTicketsAmount)
	var calcResult = float64(diference * 100)

	sorteio.Percentage = calcResult

	if err != nil {
		return models.Sorteio{}, false
	}

	return sorteio, true
}

func (s *SorteioService) GetASorteioFromObjId(sorteioObjId string) (models.Sorteio, bool) {
	var sorteio models.Sorteio

	objId, _ := primitive.ObjectIDFromHex(sorteioObjId)

	result, _ := DBConnection.FindADocument("sorteio", bson.M{"_id": objId})
	err := result.Decode(&sorteio)

	var maxTicketsAmount = sorteio.Tickets.Amount
	var ticketsSold = sorteio.TicketsSold

	var diference = float64(ticketsSold) / float64(maxTicketsAmount)
	var calcResult = float64(diference * 100)

	sorteio.Percentage = calcResult

	if err != nil {
		return models.Sorteio{}, false
	}

	return sorteio, true
}

func (s *SorteioService) GenerateAWinnerToSorteio(sorteioId string) models.Ticket {
	singleSorteio, response := s.GetASorteioFromId(sorteioId)

	if !response {
		panic(errors.New("Sorteio não encontrado"))
	}

	ticket := s.TicketService.GetWinnerTicket(singleSorteio.SorteioId, singleSorteio.Tickets.Amount)

	fmt.Println("TICKET_WINNER :: ", ticket)

	return ticket
}

func (s *SorteioService) UpdateSorteio(data dto.UpdateSorteioDto) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("sorteio")
	defer cancel()

	if validationErr := s.Validate.Struct(&data); validationErr != nil {
		panic(validationErr)
	}

	singleSorteio, result := s.GetASorteioFromId(data.SorteioId)

	if !result {
		panic(errors.New("Sorteio não encontrado"))
	}

	filter := bson.M{"_id": singleSorteio.ID}

	var sorteioStatus bool

	if data.Status == "true" {
		sorteioStatus = true
	} else {
		sorteioStatus = false
	}

	var displayFinishDateStatus bool

	if data.DisplayFinishDate == "true" {
		displayFinishDateStatus = true
	} else {
		displayFinishDateStatus = false
	}

	ticketPrice, _ := strconv.Atoi(data.TicketsPrice)
	ticketAmount, _ := strconv.Atoi(data.TicketsAmount)
	minimalForOrder, _ := strconv.Atoi(data.MinimalTicketForOrder)
	maximumForOrder, _ := strconv.Atoi(data.MaximumTicketForOrder)

	discountAmount, _ := strconv.Atoi(data.DiscountAmount)
	ticketsMinimalToDiscount, _ := strconv.Atoi(data.TicketsMinimalToDiscount)
	finishDate, _ := strconv.ParseInt(data.FinishDate, 0, 64)

	update := bson.M{
		"$set": bson.M{
			"name": data.Name,
			// "images":      data.Images,
			"description": data.Description,
			"tickets": map[string]interface{}{
				"amount":          ticketAmount,
				"price":           ticketPrice,
				"minimalForOrder": minimalForOrder,
				"maximumForOrder": maximumForOrder,
			},
			"discount": map[string]interface{}{
				"amount":  ticketsMinimalToDiscount,
				"minimal": discountAmount,
			},
			"status":              sorteioStatus,
			"finishDate":          finishDate,
			"display_finish_date": displayFinishDateStatus,
		},
	}

	_, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		panic(err)
	}
}

func (s *SorteioService) UpdateTicketsSold(sorteioId string, sold int) {
	sorteio, res := s.GetASorteioFromId(sorteioId)

	if !res {
		panic(errors.New("Sorteio não encontrado"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("sorteio")
	defer cancel()

	filter := bson.M{"_id": sorteio.ID}

	update := bson.M{
		"$set": bson.M{
			"tickets_sold": sorteio.TicketsSold + sold,
		},
	}

	_, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		panic(err)
	}
}
