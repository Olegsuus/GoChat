package handlers

import (
	handlers "github.com/Olegsuus/Auth/internal/handlers/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (h *ChatHandler) Add(c *gin.Context) {
	var dto handlers.AddChatDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	var participantIDs []primitive.ObjectID
	for _, idStr := range dto.ParticipantIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID участника"})
			return
		}
		participantIDs = append(participantIDs, id)
	}

	chat, err := h.csP.Add(c.Request.Context(), participantIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании чата"})
		return
	}

	response := gin.H{
		"chat": chat,
	}

	c.JSON(http.StatusOK, response)
}
