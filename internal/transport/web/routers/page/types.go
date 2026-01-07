package page

import (
	"context"
	"github.com/gin-gonic/gin"
)

type pageHandler interface {
	AllPagesOfBook(ctx context.Context, c *gin.Context)
	PageByNumber(ctx context.Context, c *gin.Context)
}

type pageRouter struct {
	ph pageHandler
}

func New(ph pageHandler) *pageRouter {
	return &pageRouter{ph: ph}
}
