package storage

import (
	"context"
	"fmt"
	"github.com/Olegsuus/Auth/internal/models"
	storage "github.com/Olegsuus/Auth/internal/storage/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func (s *MessageStorage) GetMany(ctx context.Context, chatID primitive.ObjectID) ([]*models.Message, error) {
	const op = "storage.GetMany"

	filter := bson.M{"chat_id": chatID}
	options := options2.Find().SetSort(bson.D{{"created_at", 1}})

	cursor, err := s.db.DataBase.Collection(storage.MessageCollection).Find(ctx, filter, options)
	if err != nil {
		log.Printf("%s: %w", op, err)
		return nil, fmt.Errorf("ошибка получения сообщений чата")
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	for cursor.Next(ctx) {
		var message models.Message
		if err = cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}
