package handlers

import (
	"context"
	"github.com/Olegsuus/GoChat/internal/handlers/dto"
	"github.com/Olegsuus/GoChat/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Register godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает нового пользователя в системе
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Param        user  body      dto.RegisterNewUserDTO  true  "Данные нового пользователя"
// @Success 	 200  "OK"
// @Failure 	 400 "Неверные данные запроса"
// @Failure 	 500  "Ошибка на сервере"
// @Router       /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var dto dto.RegisterNewUserDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	user := &models.User{
		Name:      dto.Name,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  dto.Password,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := h.hP.Add(ctx, user)
	if err != nil {
		if err.Error() == "пользователь с таким email уже существует" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с такой почтой уже существует"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации пользователя"})
		return
	}

	response := gin.H{
		"_id": id,
	}

	c.JSON(http.StatusOK, response)
}
