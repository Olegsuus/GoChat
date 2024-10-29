package handlers

import (
	"github.com/Olegsuus/Auth/internal/handlers/ws"
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

func (h *ChatHandler) ServeWS(c *gin.Context) {
	userIDStr, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизованный пользователь"})
		return
	}

	userID := userIDStr.(string)

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
