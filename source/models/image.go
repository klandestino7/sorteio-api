package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Image struct {
	ID			primitive.ObjectID 		`bson:"_id,omitempty"`
	PathUrl		string					`bson:"path_url"`
	Title		string 					`bson:"title"`
	Image		string					`bson:"image"`
}