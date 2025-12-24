package user

import (
	"context"
	"github.com/gin-gonic/gin"
)

func (ur *userRouter) UserRegRouters(ctx context.Context, gr *gin.RouterGroup) {

	Users := gr.Group("/users")
	{

		Users.POST("/create", func(c *gin.Context) { ur.uh.Create(c.Request.Context(), c) })
		Users.DELETE("/", func(c *gin.Context) { ur.uh.Delete(c.Request.Context(), c) })

		//UsersAuth := Users.Group("")
		//UsersAuth.Use(middlewares.ValidateTokenAuthorization)
		//{
		//	UsersAuth.GET("/:id", func(c *gin.Context) { ur.uh.UserById(c.Request.Context(), c) })
		//	UsersAuth.DELETE("/:id", func(c *gin.Context) { ur.uh.DeleteUser(c.Request.Context(), c) })
		//}
	}

}
