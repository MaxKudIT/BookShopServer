package server

import (
	"context"
	"github.com/gin-gonic/gin"
)

func (s *server) Create() *gin.Engine {
	router := gin.Default()

	maingr := router.Group("") //теперь принадлежит основному router
	{
		s.ur.UserRegRouters(context.TODO(), maingr) //группа maingr обновляется в основном router

	}

	return router
}
