package chat_ai

import (
	"context"
	"github.com/gin-gonic/gin"
)

type aiHandler interface {
	Ask(ctx context.Context, c *gin.Context)
}
type Router struct {
	handler aiHandler
}

func New(handler aiHandler) *Router {
	return &Router{handler: handler}
}
