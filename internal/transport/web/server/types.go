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

type pageRouter interface {
	PageRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type ubRouter interface {
	UBRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type server struct {
	ur  userRouter
	br  bookRouter
	pr  pageRouter
	ubr ubRouter
}

func New(ur userRouter, br bookRouter, pr pageRouter, ubr ubRouter) *server {
	return &server{ur: ur, br: br, pr: pr, ubr: ubr}
}
