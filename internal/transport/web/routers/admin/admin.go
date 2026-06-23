package admin

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (ar *adminRouter) AdminRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	admin := gr.Group("/admin")
	{
		admin.POST("/books", middleware.VerifyTokenMiddleware(), func(c *gin.Context) {
			ar.ah.CreateBook(c.Request.Context(), c)
		})
	}
}
