package stats

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (sr *statsRouter) StatsRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	stats := gr.Group("/stats")
	{
		stats.GET("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { sr.sh.UserStats(c.Request.Context(), c) })
	}
}
