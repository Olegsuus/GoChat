package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func (s *StorageUser) Remove(ctx context.Context, id primitive.ObjectID) error {
	const op = "storage_remove"

	filter := bson.M{"_id": id}
	result, err := s.db.UserCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("ошибка при удалении пользователя", err, "op:", op)
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}

	if result.DeletedCount == 0 {
		log.Println("Пользователь не найден", op, err)
		return fmt.Errorf("пользователь не найден")
	}

	return nil
}
