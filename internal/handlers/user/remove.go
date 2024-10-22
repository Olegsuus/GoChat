package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (h *UserHandler) Remove(c *gin.Context) {
	idStr := c.GetString("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {

	}

	err = h.hP.Remove(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении пользователя"})
		return
	}

	response := gin.H{
		"Success": true,
	}

	c.JSON(http.StatusOK, response)
}
