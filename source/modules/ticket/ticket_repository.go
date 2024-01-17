package ticket

import (
	"context"
	DBConnection "sorteio-api/source/database"
	"sorteio-api/source/models"
	"time"

	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// REPOSITORY
type ITicketRepository interface {
	Create(ticketData models.Ticket) (interface{}, bool)
	FindOne(ticketId string) models.Ticket
	FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.Ticket, int, error)
	FindOneWithFilter(filter bson.M) models.Ticket
}

type TicketRepository struct {
}

func InitTicketRepository() ITicketRepository {
	return &TicketRepository{}
}

func (ticket *TicketRepository) Create(ticketData models.Ticket) (interface{}, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("ticket")
	defer cancel()

	result, err := coll.InsertOne(ctx, ticketData)

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	if result != nil {
		return nil, false
	}

	return result.InsertedID, true
}

func (ticket *TicketRepository) FindOne(ticketId string) models.Ticket {
	objId, _ := primitive.ObjectIDFromHex(ticketId)

	filter := bson.M{"_id": objId}

	var TicketResult models.Ticket

	result, _ := DBConnection.FindADocument("ticket", filter)
	result.Decode(&TicketResult)

	return TicketResult
}

func (ticket *TicketRepository) FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.Ticket, int, error) {
	var TicketsResult []models.Ticket
	var count int

	results, err, ctx := DBConnection.FindMultipleDocuments("ticket", filter, page, pageSize)
	count = results.RemainingBatchLength()

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleTicket models.Ticket
		if err = results.Decode(&singleTicket); err != nil {
			sentry.CaptureException(err)
			panic(err)
		}

		TicketsResult = append(TicketsResult, singleTicket)
	}

	return TicketsResult, count, err
}

func (ticket *TicketRepository) FindOneWithFilter(filter bson.M) models.Ticket {
	var TicketResult models.Ticket

	result, _ := DBConnection.FindADocument("ticket", filter)
	result.Decode(&TicketResult)

	return TicketResult
}
