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

func (s *ServicesUser) CheckAuth(ctx context.Context, email, password string) (*models.User, error) {
	const op = "services.CheckAuth"

	user, err := s.sP.Get(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.l.Error("Пользователь не найден", zap.String("email", email), slog.String("op", op))
			return nil, errors.New("пользователь не найден")
		}
		s.l.Error("Ошибка при получении пользователя", zap.Error(err), slog.String("op", op))
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	if err = CheckPassword(user.Password, password); err != nil {
		s.l.Error("Неверный пароль", zap.String("email", email), slog.String("op", op))
		return nil, errors.New("неверный пароль")
	}

	s.l.Info("Успешная авторизация пользователя", zap.String("email", email))
	return user, nil
}
