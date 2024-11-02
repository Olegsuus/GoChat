package service

import (
	"context"
	"fmt"
	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"time"
)

func (s *ChatService) Add(ctx context.Context, participants []primitive.ObjectID) (*models.Chat, error) {
	const op = "service.CreateChat"

	chat := &models.Chat{
		Participants: participants,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	id, err := s.csP.Add(ctx, chat)
	if err != nil {
		s.l.Error("Ошибка при создании чата", slog.String("error", err.Error()), slog.String("op", op))
		return nil, fmt.Errorf("ошибка при создании чата: %w", err)
	}

	chat.ID = id

	s.l.Info("Чат успешно создан", slog.String("chat_id", id.Hex()), slog.String("op", op))

	return chat, nil
}
