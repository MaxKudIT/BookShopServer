package orders

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ordersHandler interface {
	Create(ctx context.Context, c *gin.Context)
}

type ordersRouter struct {
	oh ordersHandler
}

func New(oh ordersHandler) *ordersRouter {
	return &ordersRouter{oh: oh}
}
