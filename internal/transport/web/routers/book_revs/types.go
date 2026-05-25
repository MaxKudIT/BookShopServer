package book_revs

import (
	"context"

	"github.com/gin-gonic/gin"
)

type bookRevsHandler interface {
	Create(ctx context.Context, c *gin.Context)
	All(ctx context.Context, c *gin.Context)
}

type bookRevsRouter struct {
	brh bookRevsHandler
}

func New(brh bookRevsHandler) *bookRevsRouter {
	return &bookRevsRouter{brh: brh}
}
