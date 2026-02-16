package cart

import (
	"context"
	"github.com/gin-gonic/gin"
)

type cartHandler interface {
	Create(ctx context.Context, c *gin.Context)
}

type cartRouter struct {
	ch cartHandler
}

func New(ch cartHandler) *cartRouter {
	return &cartRouter{ch: ch}
}
