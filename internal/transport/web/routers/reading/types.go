package reading

import (
	"context"

	"github.com/gin-gonic/gin"
)

type readingHandler interface {
	Start(ctx context.Context, c *gin.Context)
	UpdateProgress(ctx context.Context, c *gin.Context)
	Finish(ctx context.Context, c *gin.Context)
}

type readingRouter struct {
	rh readingHandler
}

func New(rh readingHandler) *readingRouter {
	return &readingRouter{rh: rh}
}
