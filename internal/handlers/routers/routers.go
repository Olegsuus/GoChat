// internal/handlers/routers/routers.go

package routers

import (
	services "github.com/Olegsuus/Auth/internal/services/user"
	"net/http"
	"path/filepath"

	ChatHandlers "github.com/Olegsuus/Auth/internal/handlers/chat"
	messageHandlers "github.com/Olegsuus/Auth/internal/handlers/message"
	"github.com/Olegsuus/Auth/internal/handlers/middleware"
	UserHandlers "github.com/Olegsuus/Auth/internal/handlers/user"
	"github.com/Olegsuus/Auth/internal/models"
	"github.com/Olegsuus/Auth/internal/tokens/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	userHandler *UserHandlers.UserHandler,
	tokenManager *jwt.JWTManager,
	chatHandler *ChatHandlers.ChatHandler,
	messageHandler *messageHandlers.MessageHandler,
	userService *services.ServicesUser,
) *gin.Engine {
	router := gin.Default()

	router.Static("/static", filepath.Join("static"))

	router.LoadHTMLGlob("internal/templates/*.html")

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

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	router.GET("/register", func(c *gin.Context) {
		user := getUserFromContext(c)
		c.HTML(http.StatusOK, "register.html", gin.H{
			"User": user,
		})
	})

	router.GET("/login", func(c *gin.Context) {
		user := getUserFromContext(c)
		c.HTML(http.StatusOK, "login.html", gin.H{
			"User": user,
		})
	})

	router.GET("/profile", func(c *gin.Context) {
		user := getUserFromContext(c)
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"User": user,
		})
	})

	router.GET("/chats", func(c *gin.Context) {
		user := getUserFromContext(c)
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		c.HTML(http.StatusOK, "chats.html", gin.H{
			"User": user,
		})
	})

	router.GET("/chats/:id", func(c *gin.Context) {
		user := getUserFromContext(c)
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		c.HTML(http.StatusOK, "chat.html", gin.H{
			"UserID": user.ID.Hex(),
			"ChatID": c.Param("id"),
		})
	})

	router.GET("/reset-password", func(c *gin.Context) {
		user := getUserFromContext(c)
		c.HTML(http.StatusOK, "reset_password.html", gin.H{
			"User": user,
		})
	})

	router.GET("/auth/google/callback", func(c *gin.Context) {
		user := getUserFromContext(c)
		c.HTML(http.StatusOK, "google_callback.html", gin.H{
			"User": user,
		})
	})

	router.GET("/logout", func(c *gin.Context) {
		c.HTML(http.StatusOK, "logout.html", nil)
	})

	return router
}

func getUserFromContext(c *gin.Context) *models.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*models.User)
}

func getUserIDFromContext(c *gin.Context) string {
	userID, exists := c.Get("userID")
	if !exists {
		return ""
	}
	return userID.(string)
}
