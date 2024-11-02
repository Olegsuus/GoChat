package service

import (
	"context"
	"log/slog"

	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService struct {
	csP ChatStorageProvider
	l   *slog.Logger
}

type ChatStorageProvider interface {
	Add(ctx context.Context, chat *models.Chat) (primitive.ObjectID, error)
	Get(ctx context.Context, id primitive.ObjectID) (*models.Chat, error)
}

func RegisterChatService(csP ChatStorageProvider, l *slog.Logger) *ChatService {
	return &ChatService{
		csP: csP,
		l:   l,
	}
}
