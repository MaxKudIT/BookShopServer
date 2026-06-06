package chat_ai

import (
	"context"
	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) AIRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	ai := gr.Group("/ai")
	{
		ai.POST("/ask", middleware.VerifyTokenMiddleware(), func(c *gin.Context) {
			r.handler.Ask(c.Request.Context(), c)
		})
	}
}
