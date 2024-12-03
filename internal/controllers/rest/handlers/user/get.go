package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get godoc
// @Summary      Получение информации о пользователе
// @Description  Возвращает информацию о пользователе по email
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Param        email  path      string  true  "Email пользователя"
// @Success 	 200  "OK"
// @Failure 	 400 "Неверные данные запроса"
// @Failure 	 500  "Ошибка на сервере"
// @Router       /user/{email} [get]
func (h *UserHandler) Get(c *gin.Context) {
	email := c.Param("email")

	ctx := context.Background()

	user, err := h.hP.Get(ctx, email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}
