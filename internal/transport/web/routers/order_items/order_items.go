package order_items

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (oir *orderItemsRouter) OrderItemsRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	orderItems := gr.Group("/orders/:orderId/items")
	{
		orderItems.POST("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { oir.oih.Create(c.Request.Context(), c) })
		orderItems.GET("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { oir.oih.All(c.Request.Context(), c) })
	}
}
