package subscription_payments

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (spr *subscriptionPaymentsRouter) SubscriptionPaymentsRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	subscriptionPayments := gr.Group("/subscription-payments")
	{
		subscriptionPayments.POST("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { spr.sph.Create(c.Request.Context(), c) })
		subscriptionPayments.GET("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { spr.sph.All(c.Request.Context(), c) })
	}
}
