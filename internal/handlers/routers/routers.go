package routers

import (
	"github.com/Olegsuus/Auth/internal/handlers/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8001"}                     // Адрес вашего фронтенда
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"} // Включаем OPTIONS
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	router.Use(cors.New(config)) // Применение CORS middleware

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.GET("/user/:email", userHandler.Get)

	return router
}
