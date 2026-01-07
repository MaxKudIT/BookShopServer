package book

import (
	"context"
	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (br *bookRouter) BookRegRouters(ctx context.Context, gr *gin.RouterGroup) {

	Books := gr.Group("/books")
	{

		Books.GET("/all", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { br.bh.AllBooks(c.Request.Context(), c) })
		Books.GET("/my", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { br.bh.AllMyBooks(c.Request.Context(), c) })
		Books.GET("/:id", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { br.bh.BookById(c.Request.Context(), c) })

		//UsersAuth := Users.Group("")
		//UsersAuth.Use(middlewares.ValidateTokenAuthorization)
		//{
		//	UsersAuth.GET("/:id", func(c *gin.Context) { ur.uh.UserById(c.Request.Context(), c) })
		//	UsersAuth.DELETE("/:id", func(c *gin.Context) { ur.uh.DeleteUser(c.Request.Context(), c) })
		//}
	}

}
