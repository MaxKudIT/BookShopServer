package fav_items

import (
	"context"
	"github.com/gin-gonic/gin"
)

type favItemsHandler interface {
	IsInFavs(ctx context.Context, c *gin.Context)
	AllFavsItems(ctx context.Context, c *gin.Context)
	Create(ctx context.Context, c *gin.Context)
	Count(ctx context.Context, c *gin.Context)
	Delete(ctx context.Context, c *gin.Context)
}

type favItemsRouter struct {
	fih favItemsHandler
}

func New(fih favItemsHandler) *favItemsRouter {
	return &favItemsRouter{fih: fih}
}
