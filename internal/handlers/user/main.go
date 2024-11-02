package handlers

import (
	"context"
	"encoding/base64"
	"github.com/Olegsuus/GoChat/internal/tokens/jwt"
	"time"

	"crypto/rand"
	"github.com/Olegsuus/GoChat/internal/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	hP           HandlerProvider
	tokenManager *jwt.JWTManager
	oauthConfig  *oauth2.Config
}

type HandlerProvider interface {
	Add(ctx context.Context, user *models.User) (primitive.ObjectID, error)
	Get(ctx context.Context, email string) (*models.User, error)
	CheckAuth(ctx context.Context, email, password string) (*models.User, error)
	ResetPassword(ctx context.Context, email, secretWord, newPassword string) error
	UpdateProfile(ctx context.Context, id primitive.ObjectID, dto models.UpdateUserDTO) error
	HandleGoogleUser(ctx context.Context, userInfo models.GoogleUserInfo) (*models.User, error)
	Remove(ctx context.Context, id primitive.ObjectID) error
}

func RegisterHandlers(hP HandlerProvider, tokenManager *jwt.JWTManager, cfg *config.Config) *UserHandler {
	oauthConfig := &oauth2.Config{
		RedirectURL:  cfg.Google.RedirectUrl,
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		Scopes: []string{
			cfg.Google.GoogleURLEmail,
			cfg.Google.GoogleURLProfile,
		},
		Endpoint: google.Endpoint,
	}

	return &UserHandler{
		hP:           hP,
		tokenManager: tokenManager,
		oauthConfig:  oauthConfig,
	}
}

func generateStateOauthCookie(c *gin.Context) string {
	expiration := time.Now().Add(1 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", state, int(expiration.Sub(time.Now()).Seconds()), "/", "localhost", false, true)
	return state
}
