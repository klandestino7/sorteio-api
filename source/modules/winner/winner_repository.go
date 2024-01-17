package winner

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
type IWinnerRepository interface {
	Create(winnerData models.Winner) (interface{}, bool)
	FindOne(winnerId string) models.Winner
	FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.Winner, int, error)
	FindOneWithFilter(filter bson.M) models.Winner
}

type WinnerRepository struct {
}

func InitWinnerRepository() IWinnerRepository {
	return &WinnerRepository{}
}

func (winner *WinnerRepository) Create(winnerData models.Winner) (interface{}, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("winner")
	defer cancel()

	result, err := coll.InsertOne(ctx, winnerData)

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	if result != nil {
		return nil, false
	}

	return result.InsertedID, true
}

func (winner *WinnerRepository) FindOne(winnerId string) models.Winner {
	objId, _ := primitive.ObjectIDFromHex(winnerId)

	filter := bson.M{"_id": objId}

	var WinnerResult models.Winner

	result, _ := DBConnection.FindADocument("winner", filter)
	result.Decode(&WinnerResult)

	return WinnerResult
}

func (winner *WinnerRepository) FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.Winner, int, error) {
	var WinnersResult []models.Winner
	var count int

	results, err, ctx := DBConnection.FindMultipleDocuments("winner", filter, page, pageSize)
	count = results.RemainingBatchLength()

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleWinner models.Winner
		if err = results.Decode(&singleWinner); err != nil {
			sentry.CaptureException(err)
			panic(err)
		}

		WinnersResult = append(WinnersResult, singleWinner)
	}

	return WinnersResult, count, err
}

func (winner *WinnerRepository) FindOneWithFilter(filter bson.M) models.Winner {
	var WinnerResult models.Winner

	result, _ := DBConnection.FindADocument("winner", filter)
	result.Decode(&WinnerResult)

	return WinnerResult
}
