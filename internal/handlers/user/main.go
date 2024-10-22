package handlers

import (
	"context"

	"github.com/Olegsuus/Auth/internal/models"
	"github.com/Olegsuus/Auth/internal/tokens"
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
	ResetPassword(ctx context.Context, email, secretWord, newPassword string) error
	UpdateProfile(ctx context.Context, id primitive.ObjectID, dto models.UpdateUserDTO) error
	Remove(ctx context.Context, id primitive.ObjectID) error
}

func RegisterHandlers(hP HandlerProvider, tokenManager *tokens.JWTManager) *UserHandler {
	return &UserHandler{
		hP:           hP,
		tokenManager: tokenManager,
	}
}
