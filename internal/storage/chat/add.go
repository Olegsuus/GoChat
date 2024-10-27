package storage

import (
	"context"
	"fmt"
	"github.com/Olegsuus/Auth/internal/models"
	storage "github.com/Olegsuus/Auth/internal/storage/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func (s *ChatStorage) Add(ctx context.Context, chat *models.Chat) (primitive.ObjectID, error) {
	const op = "storage.AddChat"

	chat.CreatedAt = time.Now()
	chat.UpdatedAt = time.Now()

	result, err := s.db.DataBase.Collection(storage.ChatCollection).InsertOne(ctx, chat)
	if err != nil {
		log.Printf("%s: %w", op, err)
		return primitive.NilObjectID, fmt.Errorf("ошибка при создании чата: %s", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return primitive.NilObjectID, fmt.Errorf("ошибка при получении id добавленого чата: %s", err)
	}

	chat.ID = id

	log.Printf("Успешное добавление чата")

	return id, nil
}
