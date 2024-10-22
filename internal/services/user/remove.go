package services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (s *ServicesUser) Remove(ctx context.Context, id primitive.ObjectID) error {
	const op = "services.delete_user"
	s.l.With(slog.String("op", op))

	if err := s.sP.Remove(ctx, id); err != nil {
		s.l.Error("Ошибка при удалении пользователя", err.Error(), slog.String("op", op))
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}

	s.l.Info("Пользователь успешно удален", slog.String("userID", id.Hex()), slog.String("op", op))
	return nil
}
