package orders

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (or *ordersRouter) OrdersRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	orders := gr.Group("/orders")
	{
		orders.POST("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { or.oh.Create(c.Request.Context(), c) })
	}
}
