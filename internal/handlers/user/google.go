package handlers

import (
	"encoding/json"
	response2 "github.com/Olegsuus/GoChat/internal/handlers/response"
	"github.com/Olegsuus/GoChat/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GoogleCallback godoc
// @Summary      Обратный вызов после аутентификации через Google
// @Description  Обрабатывает ответ от Google после аутентификации и генерирует JWT-токен
// @Tags         Аутентификация
// @Accept       json
// @Produce      json
// @Param        code   query     	string  true  "Код авторизации от Google"
// @Param        state  query     	string  true  "Состояние запроса"
// @Success 	 200  "OK"
// @Failure 	 400 "Неверные данные запроса"
// @Failure 	 500  "Ошибка на сервере"
// @Router       /auth/google/callback [get]
func (h *UserHandler) GoogleCallback(c *gin.Context) {
	stateSaved, err := c.Cookie("oauthstate")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствует сохраненный параметр state"})
		return
	}

	stateFromProvider := c.Query("state")
	if stateFromProvider != stateSaved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный параметр state"})
		return
	}

	c.SetCookie("oauthstate", "", -1, "/", "localhost", false, true)

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствует параметр code"})
		return
	}

	token, err := h.oauthConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось обменять код на токен"})
		return
	}

	client := h.oauthConfig.Client(c.Request.Context(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении информации о пользователе"})
		return
	}
	defer resp.Body.Close()

	var googleUser models.GoogleUserInfo

	if err = json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при декодировании информации о пользователе"})
		return
	}

	user, err := h.hP.HandleGoogleUser(c.Request.Context(), googleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке пользователя"})
		return
	}

	jwtToken, err := h.tokenManager.Generate(user.ID.Hex(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	response := response2.TokenResponse{
		Message: "Аутентификация прошла успешно",
		Token:   jwtToken,
	}

	c.JSON(http.StatusOK, response)
}
