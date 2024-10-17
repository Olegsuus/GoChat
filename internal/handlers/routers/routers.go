package routers

import (
	"github.com/Olegsuus/Auth/internal/handlers/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.GET("/user/:email", userHandler.Get)

	return router
}
