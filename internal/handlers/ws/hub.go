package ws

import (
	"context"
	messageHandlers "github.com/Olegsuus/Auth/internal/handlers/message"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"sync"
)

type Hub struct {
	Clients        map[string]*Client
	Register       chan *Client
	Unregister     chan *Client
	Broadcast      chan *Message
	Mutex          sync.Mutex
	MessageService messageHandlers.MessageServiceProvider
}

func NewHub(messageService messageHandlers.MessageServiceProvider) *Hub {
	return &Hub{
		Clients:        make(map[string]*Client),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Broadcast:      make(chan *Message),
		MessageService: messageService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client.ID] = client
			h.Mutex.Unlock()
		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.Mutex.Unlock()
		case message := <-h.Broadcast:
			chatID, err := primitive.ObjectIDFromHex(message.ChatID)
			if err != nil {
				log.Printf("ошибка перевода id чата из строки в primitive.objectID: %s", err)
				continue
			}
			senderID, err := primitive.ObjectIDFromHex(message.SenderID)
			if err != nil {
				log.Printf("ошибка перевода id отправителя сообщения из строки в primitive.objectID: %s", err)
				continue
			}
			_, err = h.MessageService.SendMessage(context.Background(), chatID, senderID, message.Content)
			if err != nil {
				log.Printf("%s", err)
				continue
			}

			h.Mutex.Lock()
			for _, client := range h.Clients {
				if client.ChatID == message.ChatID {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients, client.ID)
					}
				}
			}
			h.Mutex.Unlock()
		}
	}
}
