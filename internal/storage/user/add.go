package storage

import (
	"context"
	"fmt"
	"log"

	models "Auth/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *StorageUser) Add(ctx context.Context, user *models.User) (primitive.ObjectID, error) {
	const op = "storage.Add"

	result, err := s.db.Collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return primitive.NilObjectID, fmt.Errorf("пользователь с таким email уже существует")
		}
		return primitive.NilObjectID, fmt.Errorf("ошибка при добавлении пользователя: %w", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("ошибка при получении добавленного id (%s): %w", op, err)
	}

	user.ID = id

	log.Printf("Успешное добавление нового пользователя")

	return id, nil
}
