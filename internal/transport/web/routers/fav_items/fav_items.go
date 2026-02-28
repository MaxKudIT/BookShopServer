package fav_items

import (
	"context"
	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (fir *favItemsRouter) FavItemsRegRouters(ctx context.Context, gr *gin.RouterGroup) {

	FavItems := gr.Group("/fi")
	{
		FavItems.GET("/all", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { fir.fih.AllFavsItems(c.Request.Context(), c) })
		FavItems.GET("/count", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { fir.fih.Count(c.Request.Context(), c) })
		FavItems.POST("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { fir.fih.Create(c.Request.Context(), c) })
		FavItems.POST("/infavs", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { fir.fih.IsInFavs(c.Request.Context(), c) })
		FavItems.DELETE("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { fir.fih.Delete(c.Request.Context(), c) })
	}

}
