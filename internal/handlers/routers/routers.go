package routers

import (
	ChatHandler "github.com/Olegsuus/Auth/internal/handlers/chat"
	messageHandlers "github.com/Olegsuus/Auth/internal/handlers/message"
	"github.com/Olegsuus/Auth/internal/handlers/middleware"
	handlers "github.com/Olegsuus/Auth/internal/handlers/user"
	"github.com/Olegsuus/Auth/internal/tokens"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	userHandler *handlers.UserHandler,
	tokenManager *tokens.JWTManager,
	chatHandler *ChatHandler.ChatHandler,
	messageHandler *messageHandlers.MessageHandler,
) *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	router.Use(cors.New(config))

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.GET("/user/:email", userHandler.Get)
	router.GET("/auth/google/login", userHandler.GoogleLogin)
	router.GET("/auth/google/callback", userHandler.GoogleCallback)

	authGroup := router.Group("/user")
	authGroup.Use(middleware.AuthMiddleware(tokenManager))
	{
		authGroup.POST("/password/reset", userHandler.ResetPassword)
		authGroup.PATCH("/profile", userHandler.UpdateProfile)
		authGroup.DELETE("/user", userHandler.Remove)
	}

	chatGroup := router.Group("/chats")
	chatGroup.Use(middleware.AuthMiddleware(tokenManager))
	{
		chatGroup.POST("/", chatHandler.Add)
		chatGroup.GET("/:id", chatHandler.Get)
		chatGroup.GET("/ws", chatHandler.ServeWS)
	}

	messageGroup := router.Group("/messages")
	messageGroup.Use(middleware.AuthMiddleware(tokenManager))
	{
		messageGroup.POST("/", messageHandler.SendMessage)
		messageGroup.GET("/chat/:chat_id", messageHandler.GetMessages)
	}

	return router
}
