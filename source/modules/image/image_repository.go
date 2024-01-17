package image

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
type IImageRepository interface {
	Create(ImageData models.Image) (interface{}, bool)
	FindOne(ImageId string) models.Image
}

type ImageRepository struct {
}

func InitImageRepository() IImageRepository {
	return &ImageRepository{}
}

func (image *ImageRepository) Create(ImageData models.Image) (interface{}, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("image")
	defer cancel()

	result, err := coll.InsertOne(ctx, ImageData)

	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	if result != nil {
		return nil, false
	}

	return result.InsertedID, true
}

func (image *ImageRepository) FindOne(ImageId string) models.Image {
	objId, _ := primitive.ObjectIDFromHex(ImageId)
	filter := bson.M{"_id": objId}

	var ImageResult models.Image

	result, _ := DBConnection.FindADocument("image", filter)
	result.Decode(&ImageResult)

	return ImageResult
}
