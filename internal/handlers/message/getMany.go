package handlers

import (
	resp "github.com/Olegsuus/GoChat/internal/handlers/response"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// GetMessages godoc
// @Summary      Получение сообщений из чата
// @Description  Возвращает все сообщения из указанного чата
// @Tags         Сообщения
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chat_id  path      string  true  "ID чата"
// @Success 	 200  "OK"
// @Failure 	 400 "Неверные данные запроса"
// @Failure 	 500  "Ошибка на сервере"
// @Router       /messages/chat/{chat_id} [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
	chatIDStr := c.Param("chat_id")
	chatID, err := primitive.ObjectIDFromHex(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID чата"})
		return
	}

	messages, err := h.msP.GetMessages(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении сообщений"})
		return
	}

	response := resp.MessagesResponse{
		Messages: messages,
	}

	c.JSON(http.StatusOK, response)
}
