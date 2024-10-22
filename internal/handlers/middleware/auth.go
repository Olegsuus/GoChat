package middleware

import (
	"github.com/Olegsuus/Auth/internal/tokens"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(tokenManager *tokens.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
			return
		}

		claims, err := tokenManager.Validate(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
