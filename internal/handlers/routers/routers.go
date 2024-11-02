// internal/handlers/routers/routers.go

package routers

import (
	ChatHandlers "github.com/Olegsuus/GoChat/internal/handlers/chat"
	messageHandlers "github.com/Olegsuus/GoChat/internal/handlers/message"
	"github.com/Olegsuus/GoChat/internal/handlers/middleware"
	UserHandlers "github.com/Olegsuus/GoChat/internal/handlers/user"
	services "github.com/Olegsuus/GoChat/internal/services/user"
	"github.com/Olegsuus/GoChat/internal/tokens/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(
	userHandler *UserHandlers.UserHandler,
	tokenManager *jwt.JWTManager,
	chatHandler *ChatHandlers.ChatHandler,
	messageHandler *messageHandlers.MessageHandler,
	userService *services.ServicesUser,
) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	api := router.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.GET("/user/:email", userHandler.Get)
		api.GET("/auth/google/login", userHandler.GoogleLogin)
		api.GET("/auth/google/callback", userHandler.GoogleCallback)

		authGroup := api.Group("/user")
		authGroup.Use(middleware.AuthMiddleware(tokenManager, userService))
		{
			authGroup.POST("/password/reset", userHandler.ResetPassword)
			authGroup.PATCH("/profile", userHandler.UpdateProfile)
			authGroup.DELETE("/user", userHandler.Remove)
		}

		chatGroup := api.Group("/chats")
		chatGroup.Use(middleware.AuthMiddleware(tokenManager, userService))
		{
			chatGroup.POST("/", chatHandler.Add)
			chatGroup.GET("/:id", chatHandler.Get)
			chatGroup.GET("/ws", chatHandler.ServeWS)
		}

		messageGroup := api.Group("/messages")
		messageGroup.Use(middleware.AuthMiddleware(tokenManager, userService))
		{
			messageGroup.POST("/", messageHandler.SendMessage)
			messageGroup.GET("/chat/:chat_id", messageHandler.GetMessages)
		}
	}

	return router
}
