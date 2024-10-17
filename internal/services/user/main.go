package user

import (
	"context"
	"github.com/Olegsuus/Auth/internal/models"
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
}

func RegisterServices(sP ServicesProvider, l *slog.Logger) *ServicesUser {
	return &ServicesUser{
		sP: sP,
		l:  l,
	}
}
