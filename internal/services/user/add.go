package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"log/slog"
)

func (s *ServicesUser) Add(ctx context.Context, user *models.User) (primitive.ObjectID, error) {
	const op = "services.Add"

	_, err := s.sP.Get(ctx, user.Email)
	if err == nil {
		s.l.Error("Пользователь с таким email уже существует", zap.String("email", user.Email), slog.String("op", op))
		return primitive.NilObjectID, fmt.Errorf("пользователь с таким email уже существует")
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		s.l.Error("Ошибка при проверке существования пользователя", zap.Error(err), slog.String("op", op))
		return primitive.NilObjectID, fmt.Errorf("ошибка при проверке существования пользователя: %w", err)
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		s.l.Error("Ошибка при хешировании пароля", zap.Error(err), slog.String("op", op))
		return primitive.NilObjectID, fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	user.Password = hashedPassword

	id, err := s.sP.Add(ctx, user)
	if err != nil {
		if err.Error() == "пользователь с таким email уже существует" {
			s.l.Error("Пользователь с таким email уже существует", zap.String("email", user.Email))
			return primitive.NilObjectID, err
		}
		s.l.Error("Ошибка при добавлении нового пользователя", zap.Error(err))
		return primitive.NilObjectID, fmt.Errorf("ошибка при добавлении нового пользователя: %w", err)
	}

	s.l.Info("Успешное добавление нового пользователя", zap.String("email", user.Email))

	return id, nil
}
