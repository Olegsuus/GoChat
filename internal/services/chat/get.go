package service

import (
	"context"
	"log/slog"

	"fmt"
	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *ChatService) Get(ctx context.Context, id primitive.ObjectID) (*models.Chat, error) {
	const op = "service.GetChatByID"

	chat, err := s.csP.Get(ctx, id)
	if err != nil {
		s.l.Error("Ошибка при получении чата", slog.String("error", err.Error()), slog.String("op", op))
		return nil, fmt.Errorf("ошибка при получении чата: %w", err)
	}

	s.l.Info("Чат успешно получен", slog.String("chat_id", id.Hex()), slog.String("op", op))

	return chat, nil
}
