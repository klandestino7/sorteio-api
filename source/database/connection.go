package DBConnection

import (
	"context"
	"fmt"
	"os"
	"sorteio-api/source/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBClient *mongo.Client

func StartMongoDBConnection() {
	mongoDBPassword := os.Getenv("MONGODB_PASSWORD")
	mongoDBConnection := os.Getenv("MONGODB_CONNECTION_STR")

	mongoDbConnection := fmt.Sprintf(mongoDBConnection, mongoDBPassword)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDbConnection))

	if err != nil {
		panic(err)
	}

	utils.DebugPrint("DB Connected")

	MongoDBClient = client
}

func MongoDatabase() *mongo.Database {
	return MongoDBClient.Database("sorteio")
}

// getting database collections
func GetCollection(collectionName string) *mongo.Collection {
	collection := MongoDatabase().Collection(collectionName)

	return collection
}

func FindADocument(collection string, filter bson.M) (*mongo.SingleResult, context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := MongoDatabase().Collection(collection)
	defer cancel()

	return coll.FindOne(ctx, filter), ctx
}
func FindMultipleDocuments(collection string, filter bson.M, page int, page_size int) (*mongo.Cursor, error, context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := MongoDatabase().Collection(collection)
	defer cancel()

	options := options.Find()

	if page_size != 0 {
		if page == 0 {
			page = 1
		}
		options.SetSkip(int64((page - 1) * page_size))
		options.SetLimit(int64(page_size))
	}

	results, err := coll.Find(ctx, filter, options)

	if err != nil {
		panic(err)
	}

	return results, err, ctx
}

type Change struct {
	Update    bson.M `bson:"update"`
	ReturnNew bool   `bson:"return_new"`
}

func CollectionCount(collection string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := MongoDatabase().Collection(collection)
	defer cancel()

	return coll.CountDocuments(ctx, bson.M{})
}
