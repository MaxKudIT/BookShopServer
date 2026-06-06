package server

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func (s *server) Create() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	maingr := router.Group("") //теперь принадлежит основному router
	{
		s.ur.UserRegRouters(context.TODO(), maingr) //группа maingr обновляется в основном router
		s.pr.PageRegRouters(context.TODO(), maingr)
		s.br.BookRegRouters(context.TODO(), maingr)
		s.ubr.UBRegRouters(context.TODO(), maingr)
		s.cr.CartRegRouters(context.TODO(), maingr)
		s.cir.CartItemsRegRouters(context.TODO(), maingr)
		s.fr.FavRegRouters(context.TODO(), maingr)
		s.fir.FavItemsRegRouters(context.TODO(), maingr)
		s.rr.ReadingRegRouters(context.TODO(), maingr)
		s.rsr.ReadingSessionsRegRouters(context.TODO(), maingr)
		s.brr.BookRevsRegRouters(context.TODO(), maingr)
		s.bvr.BookViewsRegRouters(context.TODO(), maingr)
		s.sr.StatsRegRouters(context.TODO(), maingr)
		s.recr.RecommendationRegRouters(context.TODO(), maingr)
		s.usr.UserSubscriptionsRegRouters(context.TODO(), maingr)
		s.spr.SubscriptionPaymentsRegRouters(context.TODO(), maingr)
		s.pbr.PhysicalBooksRegRouters(context.TODO(), maingr)
		s.or.OrdersRegRouters(context.TODO(), maingr)
		s.oir.OrderItemsRegRouters(context.TODO(), maingr)
		s.acr.AIChatRegRouters(context.TODO(), maingr)
		s.ai.AIRegRouters(context.TODO(), maingr)
	}

	return router
}
