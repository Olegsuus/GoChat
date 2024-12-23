package services

import (
	"context"
	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

type ServicesUser struct {
	sP ServicesProvider
	l  *slog.Logger
}

type ServicesProvider interface {
	Add(ctx context.Context, user *models.User) (primitive.ObjectID, error)
	Get(ctx context.Context, email string) (*models.User, error)
	CheckAuth(ctx context.Context, email, password string) (*models.User, error)
	UpdatePassword(ctx context.Context, id primitive.ObjectID, hashedPassword string) error
	UpdateProfile(ctx context.Context, id primitive.ObjectID, updateDTO bson.M) error
	Remove(ctx context.Context, id primitive.ObjectID) error
}

func RegisterServices(sP ServicesProvider, l *slog.Logger) *ServicesUser {
	return &ServicesUser{
		sP: sP,
		l:  l,
	}
}
