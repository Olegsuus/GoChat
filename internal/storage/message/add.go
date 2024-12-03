package storage

import (
	"context"
	"fmt"
	"github.com/Olegsuus/GoChat/internal/models"
	storage "github.com/Olegsuus/GoChat/internal/storage/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func (s *MessageStorage) Add(ctx context.Context, message *models.Message) (primitive.ObjectID, error) {
	const op = "storage.AddMessage"

	result, err := s.db.MessageCollection.InsertOne(ctx, message)
	if err != nil {
		log.Printf("%s: %w", op, err)
		return primitive.NilObjectID, fmt.Errorf("ошибка при добавлении сообщений в чат: %s", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("ошибка при получении id сообщения: %s", err)
	}

	log.Println("успешное добавление нового сообщения в чат")

	_, err = s.db.DataBase.Collection(storage.ChatCollection).UpdateOne(
		ctx,
		bson.M{"_id": message.ChatID},
		bson.M{
			"$set": bson.M{
				"last_message": message.ID,
				"updated_at":   time.Now(),
			},
		},
	)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("ошибка при обновлении последнего сообщения в чате")
	}

	log.Println("успешное добавление последнего сообщения в чате")

	return id, nil
}
