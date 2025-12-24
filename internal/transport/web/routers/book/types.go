package book

import (
	"context"
	"github.com/gin-gonic/gin"
)

type bookHandler interface {
	AllBooks(ctx context.Context, c *gin.Context)
	AllMyBooks(ctx context.Context, c *gin.Context)
	BookById(ctx context.Context, c *gin.Context)
}

type bookRouter struct {
	bh bookHandler
}

func New(bh bookHandler) *bookRouter {
	return &bookRouter{bh: bh}
}
