package storage

import (
	"context"
	"fmt"
	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func (s *ChatStorage) Get(ctx context.Context, id primitive.ObjectID) (*models.Chat, error) {
	const op = "storage.GetChat"

	var chat models.Chat
	err := s.db.ChatCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&chat)
	if err != nil {
		log.Printf("%s: %w", op, err)
		return nil, fmt.Errorf("ошибка при получении чата по id: %s", err)
	}

	log.Println("успешное получении чата")

	return &chat, nil
}
