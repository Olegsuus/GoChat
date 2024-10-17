package handlers

import (
	"context"

	"Auth/internal/models"
	"Auth/internal/tokens"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	hP           HandlerProvider
	tokenManager *tokens.JWTManager
}

type HandlerProvider interface {
	Add(ctx context.Context, user *models.User) (primitive.ObjectID, error)
	Get(ctx context.Context, email string) (*models.User, error)
	CheckAuth(ctx context.Context, email, password string) (*models.User, error)
}

func RegisterHandlers(hP HandlerProvider, tokenManager *tokens.JWTManager) *UserHandler {
	return &UserHandler{
		hP:           hP,
		tokenManager: tokenManager,
	}
}
