package user

import (
	"context"
	"github.com/gin-gonic/gin"
)

type userHandlers interface {
	Create(ctx context.Context, c *gin.Context)
	Delete(ctx context.Context, c *gin.Context)
}

type userRouter struct {
	uh userHandlers
}

func New(uh userHandlers) *userRouter {
	return &userRouter{uh: uh}
}
