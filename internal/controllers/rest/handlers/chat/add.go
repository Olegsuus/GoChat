package handlers

import (
	dto2 "github.com/Olegsuus/GoChat/internal/controllers/rest/handlers/dto"
	resp "github.com/Olegsuus/GoChat/internal/controllers/rest/handlers/response"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// Add godoc
// @Summary      	Создание нового чата
// @Description  	Создает новый чат с указанными участниками
// @Tags         	Чаты
// @Accept       	json
// @Produce      	json
// @Security     	BearerAuth
// @Param        	chat  body   dto.AddChatDTO  true  "ID участников чата"
// @Success 		200  "OK"
// @Failure 		400 "Неверные данные запроса"
// @Failure 		500  "Ошибка на сервере"
// @Router       	/chats [post]
func (h *ChatHandler) Add(c *gin.Context) {
	var dto dto2.AddChatDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизованный пользователь"})
		return
	}
	userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

	participantIDs := []primitive.ObjectID{userID}
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

	response := resp.ChatResponse{
		Chat: chat,
	}

	c.JSON(http.StatusOK, response)
}
