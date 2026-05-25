package book_revs

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (brr *bookRevsRouter) BookRevsRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	bookRevs := gr.Group("/book-revs")
	{
		bookRevs.POST("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { brr.brh.Create(c.Request.Context(), c) })
		bookRevs.GET("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { brr.brh.All(c.Request.Context(), c) })
	}
}
