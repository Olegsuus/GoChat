package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// Get godoc
// @Summary      Получение информации о пользователе
// @Description  Возвращает информацию о пользователе по email
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "id чата"
// @Success 	 200  "OK"
// @Failure 	 400 "Неверные данные запроса"
// @Failure 	 500  "Ошибка на сервере"
// @Router       /chats/{:id} [get]
func (h *ChatHandler) Get(c *gin.Context) {
	IDStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(IDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID чата"})
		return
	}

	chat, err := h.csP.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении чата"})
		return
	}

	response := gin.H{
		"chat": chat,
	}

	c.JSON(http.StatusOK, response)
}
