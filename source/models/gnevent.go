package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GnEvent struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Eid   string             `bson:"eid"`
	Txid  string             `bson:"txid"`
	Event string             `bson:"event"`
}

// type PixResponse struct {
// 	EndToEndId string                 `bson:"endToEndId"`
// 	Txid       string                 `bson:"txid"`
// 	Valor      string                 `bson:"valor"`
// 	Chave      string                 `bson:"chave"`
// 	Horario    string                 `bson:"horario"`
// 	Devolucoes map[string]interface{} `bson:"devolucoes"`
// }
