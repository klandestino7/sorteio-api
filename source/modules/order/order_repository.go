package order

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
type IOrderRepository interface {
	Create(orderData models.Order) (interface{}, bool)
	FindOne(orderId string) models.Order
	FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.Order, int, error)
}

type OrderRepository struct {
}

func InitOrderRepository() IOrderRepository {
	return &OrderRepository{}
}

func (order *OrderRepository) Create(orderData models.Order) (interface{}, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("Order")
	defer cancel()

	result, err := coll.InsertOne(ctx, orderData)

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	if result != nil {
		return nil, false
	}

	return result.InsertedID, true
}

func (order *OrderRepository) FindOne(orderId string) models.Order {
	objId, _ := primitive.ObjectIDFromHex(orderId)

	filter := bson.M{"_id": objId}

	var OrderResult models.Order

	result, _ := DBConnection.FindADocument("Order", filter)
	result.Decode(&OrderResult)

	return OrderResult
}

func (order *OrderRepository) FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.Order, int, error) {
	var OrdersResult []models.Order
	var count int

	results, err, ctx := DBConnection.FindMultipleDocuments("order", filter, page, pageSize)
	count = results.RemainingBatchLength()

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleOrder models.Order
		if err = results.Decode(&singleOrder); err != nil {
			sentry.CaptureException(err)
			panic(err)
		}

		OrdersResult = append(OrdersResult, singleOrder)
	}

	return OrdersResult, count, err
}
