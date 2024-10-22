package routers

import (
	"github.com/Olegsuus/Auth/internal/handlers/middleware"
	"github.com/Olegsuus/Auth/internal/handlers/user"
	"github.com/Olegsuus/Auth/internal/tokens"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(userHandler *handlers.UserHandler, tokenManager *tokens.JWTManager) *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	router.Use(cors.New(config))

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.GET("/user/:email", userHandler.Get)

	authGroup := router.Group("/user")
	authGroup.Use(middleware.AuthMiddleware(tokenManager))
	{
		authGroup.POST("/password/reset", userHandler.ResetPassword)
		//authGroup.PATCH("/profile", userHandler.)
	}
	return router
}
