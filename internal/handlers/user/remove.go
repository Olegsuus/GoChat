package handlers

import (
	resp "github.com/Olegsuus/GoChat/internal/handlers/response"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// Remove godoc
// @Summary      Удаление пользователя
// @Description  Удаляет пользователя по ID
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID пользователя"
// @Success 	 200  "OK"
// @Failure 	 400 "Неверные данные запроса"
// @Failure 	 500  "Ошибка на сервере"
// @Router       /user/{id} [delete]
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

	response := resp.SuccessResponse{
		Success: true,
	}

	c.JSON(http.StatusOK, response)
}
