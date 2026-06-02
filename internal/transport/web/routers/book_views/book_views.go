package bookviews

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (bvr *bookViewsRouter) BookViewsRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	bookViews := gr.Group("/book-views")
	{
		bookViews.POST("/:bookId", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { bvr.bvh.SaveOrUpdate(c.Request.Context(), c) })
		bookViews.GET("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { bvr.bvh.LastRecords(c.Request.Context(), c) })
	}
}
