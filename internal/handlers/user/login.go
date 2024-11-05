package handlers

import (
	"context"
	"github.com/Olegsuus/GoChat/internal/handlers/dto"
	resp "github.com/Olegsuus/GoChat/internal/handlers/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      Аутентификация пользователя
// @Description  Проверяет учетные данные пользователя и возвращает JWT-токен
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Param        credentials  body      dto.LoginDTO  true  "Учетные данные пользователя"
// @Success 	 200  "OK"
// @Failure 	 400 "Неверные данные запроса"
// @Failure 	 500  "Ошибка на сервере"
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var dto dto.LoginDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	user, err := h.hP.CheckAuth(context.Background(), dto.Email, dto.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	token, err := h.tokenManager.Generate(user.ID.Hex(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка на сервере"})
		return
	}

	response := resp.TokenResponse{
		Token: token,
	}

	c.JSON(http.StatusOK, response)
}

// GoogleLogin godoc
// @Summary      Вход через Google
// @Description  Перенаправляет пользователя на страницу авторизации Google
// @Tags         Аутентификация
// @Param        state  query     string  true  "Состояние запроса"
// @Param        code   query     string  true  "Код авторизации от Google"
// @Success      302  "Redirect to Google"
// @Failure      500  "Ошибка на сервере"
// @Router       /auth/google/login [get]
func (h *UserHandler) GoogleLogin(c *gin.Context) {
	state := generateStateOauthCookie(c)
	url := h.oauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}
