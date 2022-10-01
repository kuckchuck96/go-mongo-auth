package service

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Text      string             `json:"text"`
	Media     []byte             `json:"media"`
	Type      string             `json:"type" binding:"required,oneof=media text"`
	CreatedAt time.Time          `json:"created" bson:"createdAt"`
}

func (message *Message) New() {
	message.CreatedAt = time.Now()
}

func CreateMessage(message Message) (*mongo.InsertOneResult, error) {
	return nil, nil
}
