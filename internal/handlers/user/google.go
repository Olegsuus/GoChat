package handlers

import (
	"encoding/json"
	"github.com/Olegsuus/Auth/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

	c.JSON(http.StatusOK, gin.H{
		"message": "Аутентификация прошла успешно",
		"token":   jwtToken,
	})
}
