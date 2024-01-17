package ticket

import (
	"context"
	DBConnection "sorteio-api/source/database"
	"sorteio-api/source/models"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"

	"math/rand"
)

// SERVICE
type ITicketService interface {
	GenerateTicketNumber(sorteioId int64, numMax int) int
	CreateTicket(ticketData models.Ticket) (interface{}, error)
	CreateManyTickets(ticketData []interface{}) (interface{}, error)
	TicketHasTaken(ticketId int, sorteioId int64) bool
	CreateTicketsFromOrder()
	GetTicketFromNumber(ticketId int, sorteioId int64) models.Ticket
	GetWinnerTicket(sorteioId int64, numMax int) models.Ticket
}

type TicketService struct {
	TicketRepository ITicketRepository
	Validate         *validator.Validate
}

func InitTicketService(TicketRepository ITicketRepository, validate *validator.Validate) ITicketService {
	return &TicketService{
		TicketRepository: TicketRepository,
		Validate:         validate,
	}
}

func (s *TicketService) GenerateTicketNumber(sorteioId int64, numMax int) int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	var selectedNum int = rand.Intn(numMax-min+1) + min

	result := s.TicketHasTaken(selectedNum, sorteioId)

	if result {
		return -1
	}

	return selectedNum
}

func (s *TicketService) CreateTicket(ticketData models.Ticket) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("ticket")
	defer cancel()

	if validationErr := s.Validate.Struct(&ticketData); validationErr != nil {
		panic(validationErr)
	}

	result, err := coll.InsertOne(ctx, ticketData)

	if err != nil {
		return nil, err
	}

	if result != nil {
		return nil, err
	}

	return result.InsertedID, err
}

func (s *TicketService) CreateManyTickets(ticketData []interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("ticket")
	defer cancel()

	result, err := coll.InsertMany(ctx, ticketData)

	if err != nil {
		return nil, err
	}

	if result != nil {
		return nil, err
	}

	return result, err
}

func (s *TicketService) TicketHasTaken(ticketId int, sorteioId int64) bool {
	ticket := s.GetTicketFromNumber(ticketId, sorteioId)

	// fmt.Println("ticket :: ", ticket)

	// fmt.Println("ticket.ID.IsZero() :: ", ticket.ID.IsZero())

	if ticket.ID.IsZero() {
		return false
	}

	// fmt.Println("ticket.Number :: ", ticket.Number)
	// fmt.Println("ticketId :: ", ticketId)

	if int64(ticket.Number) == int64(ticketId) && int64(ticket.SorteioId) == int64(sorteioId) {
		return true
	}

	return false
}

func (s *TicketService) GetTicketFromNumber(ticketId int, sorteioId int64) models.Ticket {
	filter := bson.M{"number": ticketId, "sorteio_id": sorteioId}
	result, _ := DBConnection.FindADocument("ticket", filter)

	var ticket models.Ticket
	result.Decode(&ticket)

	return ticket
}

func (s *TicketService) CreateTicketsFromOrder() {

}

func (s *TicketService) GetWinnerTicket(sorteioId int64, numMax int) models.Ticket {

TRYAGAIN:
	rand.Seed(time.Now().UnixNano())
	var min = 1
	var selectedNum int = rand.Intn(numMax-min+1) + min

	ticketHasTaken := s.TicketHasTaken(selectedNum, sorteioId)

	if !ticketHasTaken {
		goto TRYAGAIN
	}

	ticket := s.GetTicketFromNumber(selectedNum, sorteioId)

	if !ticket.Status && ticket.TransactionId == "" {
		goto TRYAGAIN
	}

	return ticket
}
