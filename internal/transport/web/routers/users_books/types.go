package users_books

import (
	"context"
	"github.com/gin-gonic/gin"
)

type ubHandlers interface {
	Buy(ctx context.Context, c *gin.Context)
}

type ubRouter struct {
	ubh ubHandlers
}

func New(ubh ubHandlers) *ubRouter {
	return &ubRouter{ubh: ubh}
}
