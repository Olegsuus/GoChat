package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	LastName    string             `bson:"last_name"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	PhoneNumber string             `bson:"phone_number"`
	Country     string             `bson:"country"`
	City        string             `bson:"city"`
	SecretWord  string             `bson:"secret_word"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type UpdateUserDTO struct {
	Name        *string `json:"name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Email       *string `json:"email,omitempty"`
	Password    *string `json:"password,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Country     *string `json:"country,omitempty"`
	City        *string `json:"city,omitempty"`
}
