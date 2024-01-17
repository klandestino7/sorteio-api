package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sorteio struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	SorteioId   int64              `bson:"sorteio_id"`
	Images      []Image            `bson:"images,omitempty"`

	Tickets  Tickets  `bson:"tickets"`
	Discount Discount `bson:"discount,omitempty"`

	FinishDate        time.Time `bson:"finish_date"`
	DisplayFinishDate bool      `bson:"display_finish_date,omitempty"`

	TimeExpireOrder int `bson:"time_expire_order"`

	CreatedAt time.Time `bson:"created_at"`
	Status    bool      `bson:"status"`

	TicketsSold int     `bson:"tickets_sold, omitempty"`
	Percentage  float64 `bson:"percentage,omitempty"`
	Earning     int     `bson:"earning,omitempty"`
}

type Tickets struct {
	Amount          int `bson:"amount"`
	Price           int `bson:"price"`
	MinimalForOrder int `bson:"minimal_for_order"`
	MaximumForOrder int `bson:"maximum_for_order"`
}

type Discount struct {
	Amount  int `bson:"amount"`
	Minimal int `bson:"minimal"`
}
