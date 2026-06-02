package physical_books

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (pbr *physicalBooksRouter) PhysicalBooksRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	physicalBooks := gr.Group("/physical-books")
	{
		physicalBooks.GET("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { pbr.pbh.All(c.Request.Context(), c) })
		physicalBooks.GET("/:id", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { pbr.pbh.ById(c.Request.Context(), c) })
	}
}
