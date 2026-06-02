package subscription_payments

import (
	"context"

	"github.com/gin-gonic/gin"
)

type subscriptionPaymentsHandler interface {
	Create(ctx context.Context, c *gin.Context)
	All(ctx context.Context, c *gin.Context)
}

type subscriptionPaymentsRouter struct {
	sph subscriptionPaymentsHandler
}

func New(sph subscriptionPaymentsHandler) *subscriptionPaymentsRouter {
	return &subscriptionPaymentsRouter{sph: sph}
}
