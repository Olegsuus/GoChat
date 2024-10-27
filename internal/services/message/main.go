package service

import (
	"context"
	"log/slog"

	"github.com/Olegsuus/Auth/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceMessage struct {
	msP MessageStorageProvider
	l   *slog.Logger
}

type MessageStorageProvider interface {
	Add(ctx context.Context, message *models.Message) (primitive.ObjectID, error)
	GetMany(ctx context.Context, chatID primitive.ObjectID) ([]*models.Message, error)
}

func RegisterServiceMessage(msP MessageStorageProvider, l *slog.Logger) *ServiceMessage {
	return &ServiceMessage{
		msP: msP,
		l:   l,
	}
}
