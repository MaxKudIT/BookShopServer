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

type readingRouter interface {
	ReadingRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type readingSessionsRouter interface {
	ReadingSessionsRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type bookRevsRouter interface {
	BookRevsRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type bookViewsRouter interface {
	BookViewsRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type statsRouter interface {
	StatsRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type recommendationRouter interface {
	RecommendationRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type userSubscriptionsRouter interface {
	UserSubscriptionsRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type subscriptionPaymentsRouter interface {
	SubscriptionPaymentsRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type physicalBooksRouter interface {
	PhysicalBooksRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type ordersRouter interface {
	OrdersRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type orderItemsRouter interface {
	OrderItemsRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type aiChatRouter interface {
	AIChatRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type aiRouter interface {
	AIRegRouters(ctx context.Context, gr *gin.RouterGroup)
}

type server struct {
	ur   userRouter
	br   bookRouter
	pr   pageRouter
	ubr  ubRouter
	cir  cartItemsRouter
	cr   cartRouter
	fir  favItemsRouter
	fr   favRouter
	rr   readingRouter
	rsr  readingSessionsRouter
	brr  bookRevsRouter
	bvr  bookViewsRouter
	sr   statsRouter
	recr recommendationRouter
	usr  userSubscriptionsRouter
	spr  subscriptionPaymentsRouter
	pbr  physicalBooksRouter
	or   ordersRouter
	oir  orderItemsRouter
	acr  aiChatRouter
	ai   aiRouter
}

func New(ur userRouter, br bookRouter, pr pageRouter, ubr ubRouter, cir cartItemsRouter, cr cartRouter, fir favItemsRouter, fr favRouter, rr readingRouter, rsr readingSessionsRouter, brr bookRevsRouter, bvr bookViewsRouter, sr statsRouter, recr recommendationRouter, usr userSubscriptionsRouter, spr subscriptionPaymentsRouter, pbr physicalBooksRouter, or ordersRouter, oir orderItemsRouter, acr aiChatRouter, ai aiRouter) *server {
	return &server{ur: ur, br: br, pr: pr, ubr: ubr, cir: cir, cr: cr, fir: fir, fr: fr, rr: rr, rsr: rsr, brr: brr, bvr: bvr, sr: sr, recr: recr, usr: usr, spr: spr, pbr: pbr, or: or, oir: oir, acr: acr, ai: ai}
}
