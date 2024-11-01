package handlers

import (
	handlers "github.com/Olegsuus/Auth/internal/handlers/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (h *MessageHandler) SendMessage(c *gin.Context) {
	var dto handlers.SendMessageDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	chatID, err := primitive.ObjectIDFromHex(dto.ChatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID чата"})
		return
	}

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизованный"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	message, err := h.msP.SendMessage(c.Request.Context(), chatID, userID, dto.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при отправке сообщения"})
		return
	}

	response := gin.H{
		"message": message,
	}

	c.JSON(http.StatusOK, response)
}
