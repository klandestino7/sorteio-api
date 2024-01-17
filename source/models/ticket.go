package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Number        int                `bson:"number"`
	UserId        primitive.ObjectID `bson:"user_id"`
	SorteioId     int64              `bson:"sorteio_id"`
	OrderId       primitive.ObjectID `bson:"order_id"`
	Status        bool               `bson:"status"`
	CreatedAt     time.Time          `bson:"created_at"`
	TransactionId string             `bson:"transaction_id"`
}
