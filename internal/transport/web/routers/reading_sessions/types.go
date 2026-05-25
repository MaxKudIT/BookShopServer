package reading_sessions

import (
	"context"

	"github.com/gin-gonic/gin"
)

type readingSessionsHandler interface {
	Create(ctx context.Context, c *gin.Context)
	All(ctx context.Context, c *gin.Context)
}

type readingSessionsRouter struct {
	rsh readingSessionsHandler
}

func New(rsh readingSessionsHandler) *readingSessionsRouter {
	return &readingSessionsRouter{rsh: rsh}
}
