package user_subscriptions

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (usr *userSubscriptionsRouter) UserSubscriptionsRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	userSubscriptions := gr.Group("/user-subscriptions")
	{
		userSubscriptions.POST("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { usr.ush.Create(c.Request.Context(), c) })
		userSubscriptions.GET("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { usr.ush.All(c.Request.Context(), c) })
		userSubscriptions.GET("/status", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { usr.ush.Status(c.Request.Context(), c) })
	}
}
