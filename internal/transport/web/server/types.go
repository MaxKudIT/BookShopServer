package server

import (
	"context"
	"github.com/gin-gonic/gin"
)

type userRouter interface {
	UserRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type bookRouter interface {
	BookRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type server struct {
	ur userRouter
	br bookRouter
}

func New(ur userRouter, br bookRouter) *server {
	return &server{ur: ur, br: br}
}
