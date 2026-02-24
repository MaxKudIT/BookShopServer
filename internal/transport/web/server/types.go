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

type cartItemsRouter interface {
	CartItemsRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type cartRouter interface {
	CartRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type favItemsRouter interface {
	FavItemsRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type favRouter interface {
	FavRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type server struct {
	ur  userRouter
	br  bookRouter
	pr  pageRouter
	ubr ubRouter
	cir cartItemsRouter
	cr  cartRouter
	fir favItemsRouter
	fr  favRouter
}

func New(ur userRouter, br bookRouter, pr pageRouter, ubr ubRouter, cir cartItemsRouter, cr cartRouter, fir favItemsRouter, fr favRouter) *server {
	return &server{ur: ur, br: br, pr: pr, ubr: ubr, cir: cir, cr: cr, fir: fir, fr: fr}
}
