package user

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
type IUserRepository interface {
	Create(userData models.User) (interface{}, bool)
	FindOne(userId string) models.User
	FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.User, int, error)
	FindOneWithFilter(filter bson.M) models.User
}

type UserRepository struct {
}

func InitUserRepository() IUserRepository {
	return &UserRepository{}
}

func (user *UserRepository) Create(userData models.User) (interface{}, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("user")
	defer cancel()

	result, err := coll.InsertOne(ctx, userData)

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	if result != nil {
		return nil, false
	}

	return result.InsertedID, true
}

func (user *UserRepository) FindOne(userId string) models.User {
	objId, _ := primitive.ObjectIDFromHex(userId)

	filter := bson.M{"_id": objId}

	var UserResult models.User

	result, _ := DBConnection.FindADocument("user", filter)
	result.Decode(&UserResult)

	return UserResult
}

func (user *UserRepository) FindMultipleDocuments(filter bson.M, page int, pageSize int) ([]models.User, int, error) {
	var UsersResult []models.User
	var count int

	results, err, ctx := DBConnection.FindMultipleDocuments("user", filter, page, pageSize)
	count = results.RemainingBatchLength()

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			sentry.CaptureException(err)
			panic(err)
		}

		UsersResult = append(UsersResult, singleUser)
	}

	return UsersResult, count, err
}

func (user *UserRepository) FindOneWithFilter(filter bson.M) models.User {
	var UserResult models.User

	result, _ := DBConnection.FindADocument("user", filter)
	result.Decode(&UserResult)

	return UserResult
}
