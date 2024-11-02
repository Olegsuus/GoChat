package storage

import storage "github.com/Olegsuus/GoChat/internal/storage/mongo"

type ChatStorage struct {
	db *storage.MongoStorage
}

func RegisterStorageChat(mongoStorage *storage.MongoStorage) *ChatStorage {
	return &ChatStorage{
		db: mongoStorage,
	}
}
