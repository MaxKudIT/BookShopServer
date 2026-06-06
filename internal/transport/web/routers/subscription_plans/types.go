package subscription_plans

import (
	"context"

	"github.com/gin-gonic/gin"
)

type subscriptionPlansHandler interface {
	All(ctx context.Context, c *gin.Context)
	ById(ctx context.Context, c *gin.Context)
	ByTitle(ctx context.Context, c *gin.Context)
}

type subscriptionPlansRouter struct {
	sph subscriptionPlansHandler
}

func New(sph subscriptionPlansHandler) *subscriptionPlansRouter {
	return &subscriptionPlansRouter{sph: sph}
}
