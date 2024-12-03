package handlers

import (
	"context"
	messageHandlers "github.com/Olegsuus/GoChat/internal/controllers/rest/handlers/message"
	ws "github.com/Olegsuus/GoChat/internal/controllers/ws"
	"github.com/Olegsuus/GoChat/internal/models"
	"github.com/Olegsuus/GoChat/internal/tokens/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatHandler struct {
	csP ChatServiceProvider
	hub *ws.Hub
	tm  *jwt.JWTManager
}

type ChatServiceProvider interface {
	Add(ctx context.Context, participants []primitive.ObjectID) (*models.Chat, error)
	Get(ctx context.Context, id primitive.ObjectID) (*models.Chat, error)
}

func RegisterChatHandler(csP ChatServiceProvider, msP messageHandlers.MessageServiceProvider, tokenManager *jwt.JWTManager,
) *ChatHandler {
	hub := ws.NewHub(msP)
	go hub.Run()
	return &ChatHandler{
		csP: csP,
		hub: hub,
		tm:  tokenManager,
	}
}
