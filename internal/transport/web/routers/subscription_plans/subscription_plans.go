package subscription_plans

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (spr *subscriptionPlansRouter) SubscriptionPlansRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	subscriptionPlans := gr.Group("/subscription-plans")
	{
		subscriptionPlans.GET("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) {
			spr.sph.All(c.Request.Context(), c)
		})
		subscriptionPlans.GET("/title/:title", middleware.VerifyTokenMiddleware(), func(c *gin.Context) {
			spr.sph.ByTitle(c.Request.Context(), c)
		})
		subscriptionPlans.GET("/:id", middleware.VerifyTokenMiddleware(), func(c *gin.Context) {
			spr.sph.ById(c.Request.Context(), c)
		})
	}
}
