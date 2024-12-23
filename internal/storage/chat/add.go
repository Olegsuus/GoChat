package storage

import (
	"context"
	"fmt"
	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func (s *ChatStorage) Add(ctx context.Context, chat *models.Chat) (primitive.ObjectID, error) {
	const op = "storage.AddChat"

	chat.CreatedAt = time.Now()
	chat.UpdatedAt = time.Now()

	result, err := s.db.ChatCollection.InsertOne(ctx, chat)
	if err != nil {
		log.Printf("%s: %v", op, err)
		return primitive.NilObjectID, fmt.Errorf("ошибка при создании чата: %w", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("ошибка при получении id добавленого чата: %s", err)
	}

	chat.ID = id

	log.Printf("Успешное добавление чата")

	return id, nil
}
