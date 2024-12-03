package handlers

import (
	"github.com/Olegsuus/GoChat/internal/controllers/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWS godoc
// @Summary      Установление WebSocket соединения для чата
// @Description  Устанавливает WebSocket соединение для обмена сообщениями в реальном времени
// @Tags         Чаты
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chat_id  query     string  true  "ID чата"
// @Success 	200  "OK"
// @Failure 	400 "Неверные данные запроса"
// @Failure 	500  "Ошибка на сервере"
// @Router       /chats/ws [get]
func (h *ChatHandler) ServeWS(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}

	claims, err := h.tm.Validate(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
		return
	}

	userID := claims.UserID

	chatIDStr := c.Query("chat_id")
	if chatIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо указать chat_id"})
		return
	}

	chatID, err := primitive.ObjectIDFromHex(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный chat_id"})
		return
	}

	chat, err := h.csP.Get(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении чата"})
		return
	}

	isParticipant := false
	for _, participant := range chat.Participants {
		if participant.Hex() == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		c.JSON(http.StatusForbidden, gin.H{"error": "Вы не являетесь участником этого чата"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении соединения до WebSocket"})
		return
	}

	client := &ws.Client{
		ID:     userID,
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan *ws.Message, 256),
		ChatID: chatIDStr,
	}

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
