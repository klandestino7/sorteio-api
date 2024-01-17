package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Winner struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	UserId      primitive.ObjectID `bson:"user_id"`
	SorteioId   string             `bson:"sorteio_id"`
	SorteioName string             `bson:"sorteio_name"`
	CotaNumber  string             `bson:"cota_number"`
	Image       string             `bson:"image"`
	CreatedAt   time.Time          `bson:"created_ad"`
}
