package fav

import (
	"context"
	"github.com/gin-gonic/gin"
)

type favHandler interface {
	Create(ctx context.Context, c *gin.Context)
}

type favRouter struct {
	fh favHandler
}

func New(fh favHandler) *favRouter {
	return &favRouter{fh: fh}
}
