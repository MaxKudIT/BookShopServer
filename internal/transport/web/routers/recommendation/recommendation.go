package recommendation

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (rr *recommendationRouter) RecommendationRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	recommendation := gr.Group("/recommendations")
	{
		recommendation.GET("/home", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { rr.rh.HomeRecommendation(c.Request.Context(), c) })
		recommendation.GET("/cart", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { rr.rh.CartRecommendation(c.Request.Context(), c) })
		recommendation.GET("/page", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { rr.rh.RecommesPage(c.Request.Context(), c) })
		recommendation.GET("/books/:bookId", middleware.VerifyTokenMiddleware(), func(c *gin.Context) { rr.rh.RecommendationByBook(c.Request.Context(), c) })
	}
}
