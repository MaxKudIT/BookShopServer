package cart_items

import (
	"context"
	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (cir *cartItemsRouter) CartItemsRegRouters(ctx context.Context, gr *gin.RouterGroup) {

	CartItems := gr.Group("/ci")
	{
		CartItems.GET("/all", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { cir.cih.AllCartItems(c.Request.Context(), c) })
		CartItems.GET("/count", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { cir.cih.Count(c.Request.Context(), c) })
		CartItems.POST("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { cir.cih.Create(c.Request.Context(), c) })
		CartItems.POST("/some", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { cir.cih.CreateItems(c.Request.Context(), c) })
		CartItems.POST("/incart", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { cir.cih.IsInCart(c.Request.Context(), c) })
		CartItems.POST("/allincart", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { cir.cih.AreAllInCart(c.Request.Context(), c) })
		CartItems.DELETE("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { cir.cih.Delete(c.Request.Context(), c) })
	}

}
