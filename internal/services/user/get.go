package user

import (
	"Auth/internal/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"log/slog"
)

func (s *ServicesUser) Get(ctx context.Context, email string) (*models.User, error) {
	const op = "services.Get"

	user, err := s.sP.Get(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.l.Error("Пользователь не найден", zap.String("email", email), slog.String("op", op))
			return nil, fmt.Errorf("не найден пользователь с такой почтой: %w", err)
		}
		s.l.Error("Ошибка при получении пользователя", zap.Error(err), slog.String("op", op))
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	s.l.Info("Пользователь успешно получен", zap.String("email", email), slog.String("op", op))

	return user, nil
}
