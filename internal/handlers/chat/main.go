package handlers

import (
	"context"
	"github.com/Olegsuus/Auth/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatHandler struct {
	csP ChatServiceProvider
}

type ChatServiceProvider interface {
	Add(ctx context.Context, participants []primitive.ObjectID) (*models.Chat, error)
	Get(ctx context.Context, id primitive.ObjectID) (*models.Chat, error)
}

func RegisterChatHandler(csP ChatServiceProvider) *ChatHandler {
	return &ChatHandler{
		csP: csP,
	}
}
