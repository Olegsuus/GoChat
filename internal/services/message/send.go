package service

import (
	"context"
	"log/slog"

	"fmt"
	"time"

	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *ServiceMessage) SendMessage(ctx context.Context, chatID, senderID primitive.ObjectID, content string) (*models.Message, error) {
	const op = "service.SendMessage"

	message := &models.Message{
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	id, err := s.msP.Add(ctx, message)
	if err != nil {
		s.l.Error("Ошибка при отправке сообщения", slog.String("error", err.Error()), slog.String("op", op))
		return nil, fmt.Errorf("ошибка при отправке сообщения: %w", err)
	}

	message.ID = id

	s.l.Info("Сообщение успешно отправлено", slog.String("message_id", id.Hex()), slog.String("op", op))

	return message, nil
}
