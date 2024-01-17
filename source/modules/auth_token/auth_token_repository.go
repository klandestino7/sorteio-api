package authToken

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
type IAuthTokenRepository interface {
	Create(authTokenData models.Token) (interface{}, bool)
	FindOne(authTokenId string) models.Token
}

type AuthTokenRepository struct {
}

func InitAuthTokenRepository() IAuthTokenRepository {
	return &AuthTokenRepository{}
}

func (authToken *AuthTokenRepository) Create(authTokenData models.Token) (interface{}, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("authToken")
	defer cancel()

	result, err := coll.InsertOne(ctx, authTokenData)

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	if result != nil {
		return nil, false
	}

	return result.InsertedID, true
}

func (authToken *AuthTokenRepository) FindOne(authTokenId string) models.Token {
	objId, _ := primitive.ObjectIDFromHex(authTokenId)
	filter := bson.M{"_id": objId}

	var authTokenResult models.Token

	result, _ := DBConnection.FindADocument("authToken", filter)
	result.Decode(&authTokenResult)

	return authTokenResult
}
