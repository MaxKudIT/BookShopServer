package users_books

import (
	"context"
	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (ubr *ubRouter) UBRegRouters(ctx context.Context, gr *gin.RouterGroup) {

	ub := gr.Group("/ub")
	{

		ub.POST("/buy", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { ubr.ubh.Buy(c.Request.Context(), c) })
	}

}
