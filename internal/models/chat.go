package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Chat struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	Participants []primitive.ObjectID `bson:"participants"`
	CreatedAt    time.Time            `bson:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at"`
	LastMessage  *primitive.ObjectID  `bson:"last_message,omitempty" json:"last_message,omitempty"`
}
