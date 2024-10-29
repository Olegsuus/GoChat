package storage

import (
	"context"
	"fmt"
	"github.com/Olegsuus/Auth/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	UserCollection    = "users"
	ChatCollection    = "chats"
	MessageCollection = "messages"
)

type MongoStorage struct {
	Client            *mongo.Client
	DataBase          *mongo.Database
	UserCollection    *mongo.Collection
	ChatCollection    *mongo.Collection
	MessageCollection *mongo.Collection
}

func NewMongoStorage(cfg *config.Config) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOprions := options.Client().ApplyURI(cfg.Mongo.URI)
	client, err := mongo.Connect(ctx, clientOprions)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MongoDB: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения с MongoDB: %w", err)
	}

	db := client.Database(cfg.Mongo.DBNAME)

	userCollection := db.Collection(UserCollection)
	chatCollection := db.Collection(ChatCollection)
	messageCollection := db.Collection(MessageCollection)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true).SetName("unique_email"),
	}

	_, err = userCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, fmt.Errorf("ошибка создаия индекса на email: %w", err)
	}

	log.Println("Подключение к MongoDB установлено")

	return &MongoStorage{
		Client:            client,
		DataBase:          db,
		UserCollection:    userCollection,
		ChatCollection:    chatCollection,
		MessageCollection: messageCollection,
	}, nil
}

func (s *MongoStorage) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.Client.Disconnect(ctx)
}
