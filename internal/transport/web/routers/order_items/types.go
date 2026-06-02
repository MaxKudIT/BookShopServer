package order_items

import (
	"context"

	"github.com/gin-gonic/gin"
)

type orderItemsHandler interface {
	Create(ctx context.Context, c *gin.Context)
	All(ctx context.Context, c *gin.Context)
}

type orderItemsRouter struct {
	oih orderItemsHandler
}

func New(oih orderItemsHandler) *orderItemsRouter {
	return &orderItemsRouter{oih: oih}
}
