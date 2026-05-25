package stats

import (
	"context"

	"github.com/gin-gonic/gin"
)

type statsHandler interface {
	UserStats(ctx context.Context, c *gin.Context)
}

type statsRouter struct {
	sh statsHandler
}

func New(sh statsHandler) *statsRouter {
	return &statsRouter{sh: sh}
}
