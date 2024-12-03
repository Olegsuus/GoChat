package middleware

import (
	"context"
	services "github.com/Olegsuus/GoChat/internal/services/user"
	"github.com/Olegsuus/GoChat/internal/tokens/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(tokenManager *jwt.JWTManager, userService *services.ServicesUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
			return
		}

		// Разделяем заголовок на "Bearer" и сам токен
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			return
		}

		tokenString := parts[1]

		claims, err := tokenManager.Validate(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			return
		}

		user, err := userService.Get(context.Background(), claims.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
			return
		}

		// Добавление данных пользователя в контекст
		c.Set("userID", user.ID.Hex())
		c.Set("user", user)

		c.Next()
	}
}
