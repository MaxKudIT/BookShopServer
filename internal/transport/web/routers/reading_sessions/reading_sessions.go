package reading_sessions

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (rsr *readingSessionsRouter) ReadingSessionsRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	readingSessions := gr.Group("/reading-sessions")
	{
		readingSessions.POST("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { rsr.rsh.Create(c.Request.Context(), c) })
		readingSessions.GET("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { rsr.rsh.All(c.Request.Context(), c) })
	}
}
