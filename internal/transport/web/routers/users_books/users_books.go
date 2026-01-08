package users_books

import (
	"context"
	"github.com/gin-gonic/gin"
)

func (ubr *ubRouter) UBRegRouters(ctx context.Context, gr *gin.RouterGroup) {

	ub := gr.Group("/ub")
	{

		ub.POST("/buy", func(c *gin.Context) { ubr.ubh.Buy(c.Request.Context(), c) })
	}

}
