package session

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
type ISessionRepository interface {
	Create(sessionData models.Session) (interface{}, bool)
	FindOne(sessionId string) models.Session
	FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.Session, int, error)
}

type SessionRepository struct {
}

func InitSessionRepository() ISessionRepository {
	return &SessionRepository{}
}

func (session *SessionRepository) Create(sessionData models.Session) (interface{}, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("session")
	defer cancel()

	result, err := coll.InsertOne(ctx, sessionData)

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	if result != nil {
		return nil, false
	}

	return result.InsertedID, true
}

func (session *SessionRepository) FindOne(sessionId string) models.Session {
	objId, _ := primitive.ObjectIDFromHex(sessionId)

	filter := bson.M{"_id": objId}

	var SessionResult models.Session

	result, _ := DBConnection.FindADocument("session", filter)
	result.Decode(&SessionResult)

	return SessionResult
}

func (session *SessionRepository) FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.Session, int, error) {
	var SessionsResult []models.Session
	var count int

	results, err, ctx := DBConnection.FindMultipleDocuments("session", filter, page, pageSize)
	count = results.RemainingBatchLength()

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleSession models.Session
		if err = results.Decode(&singleSession); err != nil {
			sentry.CaptureException(err)
			panic(err)
		}

		SessionsResult = append(SessionsResult, singleSession)
	}

	return SessionsResult, count, err
}
