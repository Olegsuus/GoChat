package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Get(c *gin.Context) {
	email := c.Param("email")

	ctx := context.Background()

	user, err := h.hP.Get(ctx, email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}
