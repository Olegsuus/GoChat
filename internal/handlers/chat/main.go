package handlers

import (
	"context"
	messageHandlers "github.com/Olegsuus/GoChat/internal/handlers/message"
	ws "github.com/Olegsuus/GoChat/internal/handlers/ws"
	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatHandler struct {
	csP ChatServiceProvider
	hub *ws.Hub
}

type ChatServiceProvider interface {
	Add(ctx context.Context, participants []primitive.ObjectID) (*models.Chat, error)
	Get(ctx context.Context, id primitive.ObjectID) (*models.Chat, error)
}

func RegisterChatHandler(csP ChatServiceProvider, msP messageHandlers.MessageServiceProvider) *ChatHandler {
	hub := ws.NewHub(msP)
	go hub.Run()
	return &ChatHandler{
		csP: csP,
		hub: hub,
	}
}
