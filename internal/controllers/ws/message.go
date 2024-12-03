package ws

type Message struct {
	SenderID string `json:"sender_id"`
	ChatID   string `json:"chat_id"`
	Content  string `json:"content"`
}
