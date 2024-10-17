package storage

import (
	storage "github.com/Olegsuus/Auth/internal/storage/mongo"
)

type StorageUser struct {
	db *storage.MongoStorage
}

func RegisterStorage(mongoStorage *storage.MongoStorage) *StorageUser {
	return &StorageUser{
		db: mongoStorage,
	}
}
