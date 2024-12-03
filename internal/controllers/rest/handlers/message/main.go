package handlers

import (
	"context"
	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageHandler struct {
	msP MessageServiceProvider
}

type MessageServiceProvider interface {
	SendMessage(ctx context.Context, chatID, senderID primitive.ObjectID, content string) (*models.Message, error)
	GetMessages(ctx context.Context, chatID primitive.ObjectID) ([]*models.Message, error)
}

func RegisterMessageHandlers(msP MessageServiceProvider) *MessageHandler {
	return &MessageHandler{
		msP: msP,
	}
}
