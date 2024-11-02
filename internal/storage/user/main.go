package storage

import (
	storage "github.com/Olegsuus/GoChat/internal/storage/mongo"
)

type StorageUser struct {
	db *storage.MongoStorage
}

func RegisterStorage(mongoStorage *storage.MongoStorage) *StorageUser {
	return &StorageUser{
		db: mongoStorage,
	}
}
