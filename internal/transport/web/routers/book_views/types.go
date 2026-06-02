package bookviews

import (
	"context"

	"github.com/gin-gonic/gin"
)

type bookViewsHandler interface {
	SaveOrUpdate(ctx context.Context, c *gin.Context)
	LastRecords(ctx context.Context, c *gin.Context)
}

type bookViewsRouter struct {
	bvh bookViewsHandler
}

func New(bvh bookViewsHandler) *bookViewsRouter {
	return &bookViewsRouter{bvh: bvh}
}
