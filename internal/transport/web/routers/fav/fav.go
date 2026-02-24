package fav

import "C"
import (
	"context"
	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (fr *favRouter) FavRegRouters(ctx context.Context, gr *gin.RouterGroup) {

	Fav := gr.Group("/fav")
	{
		Fav.POST("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { fr.fh.Create(c.Request.Context(), c) })

	}

}
