package cart

import (
	"context"
	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (cr *cartRouter) CartRegRouters(ctx context.Context, gr *gin.RouterGroup) {

	Cart := gr.Group("/cart")
	{
		Cart.POST("/", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { cr.ch.Create(c.Request.Context(), c) })

	}

}
