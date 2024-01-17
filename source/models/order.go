package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	OrderId       int64              `bson:"order_id"`
	UserId        primitive.ObjectID `bson:"user_id"`
	UserName      string             `bson:"user_name,omitempty"`
	SorteioId     int64              `bson:"sorteio_id"`
	TransactionId string             `bson:"transaction_id"`

	Tickets       []int  `bson:"tickets"`
	TicketsAmount int    `bson:"tickets_amount"`
	Status        string `bson:"status"`
	PaymentMethod string `bson:"payment_method"`
	QRCodeImage   string `bson:"qr_code_image"`
	QRCodeString  string `bson:"qr_code_string"`
	Total         int    `bson:"total"`
	Referal       string `bson:"referal,omitempty"`
	Rejected      bool   `bson:"rejected,omitempty"`

	ExpireAt  time.Time `bson:"expire_at"`
	CreatedAt time.Time `bson:"created_at"`
}

type OrderStatus struct {
	ID        int   `bson:"_id,omitempty"`
	Status    uint8 `bson:"status"`
	CreatedAt int   `bson:"created_at"`
}

// func (order *Order) GetOrderStatusAsString() string {
// 	switch order.Status {
// 	case 0:
// 		return "created"
// 	case 1:
// 		return "approved"
// 	case 2:
// 		return "rejected"
// 	default:
// 		return "unknow"
// 	}
// }
