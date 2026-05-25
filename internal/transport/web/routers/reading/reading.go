package reading

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (rr *readingRouter) ReadingRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	reading := gr.Group("/reading")
	{
		reading.POST("/start", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { rr.rh.Start(c.Request.Context(), c) })
		reading.PATCH("/progress", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { rr.rh.UpdateProgress(c.Request.Context(), c) })
		reading.POST("/finish", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { rr.rh.Finish(c.Request.Context(), c) })
	}
}
