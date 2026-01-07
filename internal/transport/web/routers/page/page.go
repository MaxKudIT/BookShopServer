package page

import (
	"context"
	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (pr *pageRouter) PageRegRouters(ctx context.Context, gr *gin.RouterGroup) {

	Books := gr.Group("/books")
	{
		Books.GET("/:id/pagesCount", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { pr.ph.AllPagesOfBook(c.Request.Context(), c) })
		Books.GET("/:id/pages/:pageNumber", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { pr.ph.PageByNumber(c.Request.Context(), c) })

	}

}
