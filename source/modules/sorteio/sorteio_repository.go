package sorteio

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
type ISorteioRepository interface {
	Create(sorteioData models.Sorteio) (interface{}, bool)
	FindOne(sorteioId string) models.Sorteio
	FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.Sorteio, int, error)
}

type SorteioRepository struct {
}

func InitSorteioRepository() ISorteioRepository {
	return &SorteioRepository{}
}

func (sorteio *SorteioRepository) Create(sorteioData models.Sorteio) (interface{}, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("Sorteio")
	defer cancel()

	result, err := coll.InsertOne(ctx, sorteioData)

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	if result != nil {
		return nil, false
	}

	return result.InsertedID, true
}

func (sorteio *SorteioRepository) FindOne(sorteioId string) models.Sorteio {
	objId, _ := primitive.ObjectIDFromHex(sorteioId)

	filter := bson.M{"_id": objId}

	var SorteioResult models.Sorteio

	result, _ := DBConnection.FindADocument("Sorteio", filter)
	result.Decode(&SorteioResult)

	return SorteioResult
}

func (sorteio *SorteioRepository) FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.Sorteio, int, error) {
	var SorteiosResult []models.Sorteio
	var count int

	results, err, ctx := DBConnection.FindMultipleDocuments("sorteio", filter, page, pageSize)
	count = results.RemainingBatchLength()

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleSorteio models.Sorteio
		if err = results.Decode(&singleSorteio); err != nil {
			sentry.CaptureException(err)
			panic(err)
		}

		SorteiosResult = append(SorteiosResult, singleSorteio)
	}

	return SorteiosResult, count, err
}
