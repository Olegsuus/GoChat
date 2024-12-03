// internal/controllers/routers/routers.go

package routers

import (
	ChatHandlers "github.com/Olegsuus/GoChat/internal/controllers/rest/handlers/chat"
	messageHandlers "github.com/Olegsuus/GoChat/internal/controllers/rest/handlers/message"
	UserHandlers "github.com/Olegsuus/GoChat/internal/controllers/rest/handlers/user"
	"github.com/Olegsuus/GoChat/internal/controllers/rest/middleware"
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

	authMiddleware := middleware.AuthMiddleware(tokenManager, userService)

	api := router.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.GET("/auth/google/login", userHandler.GoogleLogin)
		api.GET("/auth/google/callback", userHandler.GoogleCallback)

		api.GET("/user/:email", authMiddleware, userHandler.Get)

		authGroup := api.Group("/user")
		authGroup.Use(authMiddleware)
		{
			authGroup.POST("/password/reset", userHandler.ResetPassword)
			authGroup.PATCH("/profile", userHandler.UpdateProfile)
			authGroup.DELETE("/user", userHandler.Remove)
		}

		chatGroup := api.Group("/chats")
		chatGroup.Use(authMiddleware)
		{
			chatGroup.POST("/", chatHandler.Add)
			chatGroup.GET("/:id", chatHandler.Get)
			chatGroup.GET("/ws", chatHandler.ServeWS)
		}

		messageGroup := api.Group("/messages")
		messageGroup.Use(authMiddleware)
		{
			messageGroup.POST("/", messageHandler.SendMessage)
			messageGroup.GET("/chat/:chat_id", messageHandler.GetMessages)
		}
	}

	return router
}
