package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func (s *StorageUser) UpdatePassword(ctx context.Context, id primitive.ObjectID, hashedPassword string) error {
	const op = "storage_update_password"

	filter := bson.M{
		"_id": id,
	}
	update := bson.M{
		"$set": bson.M{
			"password":   hashedPassword,
			"updated_at": time.Now(),
		},
	}

	_, err := s.db.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("%s: %w", op, err)
		return fmt.Errorf("ошибка при обновлении пароля: %w", err)
	}

	log.Println("Успешное обновление пароля")

	return nil
}

func (s *StorageUser) UpdateProfile(ctx context.Context, id primitive.ObjectID, updateDTO bson.M) error {
	const op = "storage_update_profile"

	filter := bson.M{
		"_id": id,
	}

	update := bson.M{
		"$set": updateDTO,
	}

	result, err := s.db.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("%s: %w", op, err)
		return fmt.Errorf("ошибка при обновлении профиля пользователя: %w", err)
	}

	if result.MatchedCount == 0 {
		log.Printf("%s: %w", op, err)
		return fmt.Errorf("пользователь не найден")
	}
	return nil
}
