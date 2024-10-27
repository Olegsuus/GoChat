package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Olegsuus/Auth/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *ServiceMessage) GetMessages(ctx context.Context, chatID primitive.ObjectID) ([]*models.Message, error) {
	const op = "service.GetMessages"

	messages, err := s.msP.GetMany(ctx, chatID)
	if err != nil {
		s.l.Error("Ошибка при получении сообщений", slog.String("error", err.Error()), slog.String("op", op))
		return nil, fmt.Errorf("ошибка при получении сообщений: %w", err)
	}

	s.l.Info("Сообщения успешно получены", slog.String("chat_id", chatID.Hex()), slog.String("op", op))

	return messages, nil
}
