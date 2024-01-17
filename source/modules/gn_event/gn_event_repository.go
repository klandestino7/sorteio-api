package gnEvent

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
type IGnEventRepository interface {
	Create(GnEventData models.GnEvent) (interface{}, bool)
	FindOne(GnEventId string) models.GnEvent
}

type GnEventRepository struct {
}

func InitGnEventRepository() IGnEventRepository {
	return &GnEventRepository{}
}

func (gnEvent *GnEventRepository) Create(GnEventData models.GnEvent) (interface{}, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("gnEvent")
	defer cancel()

	result, err := coll.InsertOne(ctx, GnEventData)

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	if result != nil {
		return nil, false
	}

	return result.InsertedID, true
}

func (gnEvent *GnEventRepository) FindOne(GnEventId string) models.GnEvent {
	objId, _ := primitive.ObjectIDFromHex(GnEventId)
	filter := bson.M{"_id": objId}

	var GnEventResult models.GnEvent

	result, _ := DBConnection.FindADocument("gnEvent", filter)
	result.Decode(&GnEventResult)

	return GnEventResult
}
