package handlers

import (
	"context"
	handlers "github.com/Olegsuus/Auth/internal/handlers/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Login(c *gin.Context) {
	var dto handlers.LoginDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	user, err := h.hP.CheckAuth(context.Background(), dto.Email, dto.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	token, err := h.tokenManager.Generate(user.ID.Hex(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка на сервере"})
		return
	}

	response := gin.H{
		"token": token,
	}

	c.JSON(http.StatusOK, response)
}
