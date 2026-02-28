package cart_items

import (
	"context"
	"github.com/gin-gonic/gin"
)

type cartItemsHandler interface {
	IsInCart(ctx context.Context, c *gin.Context)
	AllCartItems(ctx context.Context, c *gin.Context)
	AreAllInCart(ctx context.Context, c *gin.Context)
	Create(ctx context.Context, c *gin.Context)
	CreateItems(ctx context.Context, c *gin.Context)
	Delete(ctx context.Context, c *gin.Context)
	Count(ctx context.Context, c *gin.Context)
}

type cartItemsRouter struct {
	cih cartItemsHandler
}

func New(cih cartItemsHandler) *cartItemsRouter {
	return &cartItemsRouter{cih: cih}
}
