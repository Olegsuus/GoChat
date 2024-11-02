package handlers

import (
	"github.com/Olegsuus/GoChat/internal/handlers/dto"
	resp "github.com/Olegsuus/GoChat/internal/handlers/response"
	"github.com/Olegsuus/GoChat/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// ResetPassword godoc
// @Summary      Сброс пароля пользователя
// @Description  Позволяет пользователю сбросить пароль, используя секретное слово
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        reset  body      dto.ResetPasswordDTO  true  "Данные для сброса пароля"
// @Success 	 200  "OK"
// @Failure 	 400 "Неверные данные запроса"
// @Failure 	 500  "Ошибка на сервере"
// @Router       /user/password/reset [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var dto dto.ResetPasswordDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	err := h.hP.ResetPassword(c.Request.Context(), dto.Email, dto.SecretWord, dto.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := resp.SuccessResponse{
		Success: true,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProfile godoc
// @Summary      Обновление профиля пользователя
// @Description  Обновляет данные профиля пользователя
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user  body      models.UpdateUserDTO  true  "Данные для обновления профиля"
// / @Success 	 200  "OK"
// @Failure 	 400 "Неверные данные запроса"
// @Failure 	 500  "Ошибка на сервере"
// @Router       /user/profile [patch]
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

	response := resp.SuccessResponse{
		Success: true,
	}

	c.JSON(http.StatusOK, response)
}
