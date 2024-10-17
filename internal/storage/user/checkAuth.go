package storage

import (
	"Auth/internal/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *StorageUser) CheckAuth(ctx context.Context, email, password string) (*models.User, error) {
	const op = "storage.CheckAuth"

	user, err := s.Get(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("пользователь с почтой %s не найден", email)
		}
		return nil, fmt.Errorf("ошибка при проверке пользователя: %w", err)
	}

	if user.Password != password {
		return nil, errors.New("неверный пароль")
	}

	return user, nil
}
