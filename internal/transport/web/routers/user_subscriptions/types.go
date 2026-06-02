package user_subscriptions

import (
	"context"

	"github.com/gin-gonic/gin"
)

type userSubscriptionsHandler interface {
	Create(ctx context.Context, c *gin.Context)
	All(ctx context.Context, c *gin.Context)
	Status(ctx context.Context, c *gin.Context)
}

type userSubscriptionsRouter struct {
	ush userSubscriptionsHandler
}

func New(ush userSubscriptionsHandler) *userSubscriptionsRouter {
	return &userSubscriptionsRouter{ush: ush}
}
