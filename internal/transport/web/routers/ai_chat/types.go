package ai_chat

import (
	"context"

	"github.com/gin-gonic/gin"
)

type aiChatHandler interface {
	CreateChat(ctx context.Context, c *gin.Context)
	CreateMessage(ctx context.Context, c *gin.Context)
	CurrentChat(ctx context.Context, c *gin.Context)
	Messages(ctx context.Context, c *gin.Context)
	DeleteMessages(ctx context.Context, c *gin.Context)
	Ask(ctx context.Context, c *gin.Context)
}

type aiChatRouter struct {
	ach aiChatHandler
}

func New(ach aiChatHandler) *aiChatRouter {
	return &aiChatRouter{ach: ach}
}
