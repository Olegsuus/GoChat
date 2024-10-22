package handlers

import (
	handlers "github.com/Olegsuus/Auth/internal/handlers/dto"
	"github.com/Olegsuus/Auth/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var dto handlers.ResetPasswordDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	err := h.hP.ResetPassword(c.Request.Context(), dto.Email, dto.SecretWord, dto.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"Success": true,
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	idStr := c.GetString("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный идентификатор пользователя"})
		return
	}

	var dto models.UpdateUserDTO
	if err = c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	err = h.hP.UpdateProfile(c.Request.Context(), id, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении данных пользователя"})
		return
	}

	response := gin.H{
		"Success": true,
	}

	c.JSON(http.StatusOK, response)
}
