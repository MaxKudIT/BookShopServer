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
	}

	return router
}
