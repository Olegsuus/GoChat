package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *StorageUser) Get(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.Get"

	filter := bson.M{"email": email}

	var user models.User
	err := s.db.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("не найден пользователь с такой почтой: %w", err)
		}
		return nil, fmt.Errorf("ошибка при получении пользователя по почте (%s): %w", op, err)
	}

	return &user, nil
}
