package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	ChatID    primitive.ObjectID   `bson:"chat_id"`
	SenderID  primitive.ObjectID   `bson:"sender_id"`
	Content   string               `bson:"content"`
	CreatedAt time.Time            `bson:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at"`
	ReadBy    []primitive.ObjectID `bson:"read_by,omitempty" json:"read_by,omitempty"`
}
