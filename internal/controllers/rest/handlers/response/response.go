package response

import "github.com/Olegsuus/GoChat/internal/models"

type SuccessResponse struct {
	Success bool `json:"success"`
}

type RegisterResponse struct {
	ID string `json:"_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}

type ChatResponse struct {
	Chat *models.Chat `json:"chat"`
}

type MessageResponse struct {
	Message *models.Message `json:"message"`
}

type MessagesResponse struct {
	Messages []*models.Message `json:"messages"`
}
