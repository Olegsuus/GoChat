package storage

import storage "github.com/Olegsuus/Auth/internal/storage/mongo"

type ChatStorage struct {
	db *storage.MongoStorage
}

func RegisterStorageChat(mongoStorage *storage.MongoStorage) *ChatStorage {
	return &ChatStorage{
		db: mongoStorage,
	}
}
