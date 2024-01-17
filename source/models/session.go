package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Session struct {
	ID   		primitive.ObjectID 		`bson:"_id,omitempty"`
	Token		string					`bson:"token"`
	Device		string					`bson:"device"`
	CurrentIp	string					`bson:"current_ip"`
	LastIp		string					`bson:"last_ip,omitempty"`
	UserId		primitive.ObjectID		`bson:"user_id,omitempty"`
	GeoLoc		string					`bson:"geo_loc,omitempty"`
	UpdatedAt	time.Time				`bson:"updated_at"`
}

func (session *Session) UpdateIp(newIp string) bool {
	session.LastIp = session.CurrentIp
	session.CurrentIp = newIp
	return true
}