package admin

import (
	"context"

	"github.com/gin-gonic/gin"
)

type adminHandler interface {
	CreateBook(ctx context.Context, c *gin.Context)
}

type adminRouter struct {
	ah adminHandler
}

func New(ah adminHandler) *adminRouter {
	return &adminRouter{ah: ah}
}
