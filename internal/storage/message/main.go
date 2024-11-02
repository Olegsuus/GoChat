package storage

import (
	storage "github.com/Olegsuus/GoChat/internal/storage/mongo"
)

type MessageStorage struct {
	db *storage.MongoStorage
}

func RegisterStorageMessage(mongoStorage *storage.MongoStorage) *MessageStorage {
	return &MessageStorage{
		db: mongoStorage,
	}
}
