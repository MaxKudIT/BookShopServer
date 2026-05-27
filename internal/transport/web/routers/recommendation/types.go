package recommendation

import (
	"context"

	"github.com/gin-gonic/gin"
)

type recommendationHandler interface {
	HomeRecommendation(ctx context.Context, c *gin.Context)
	CartRecommendation(ctx context.Context, c *gin.Context)
	RecommesPage(ctx context.Context, c *gin.Context)
	RecommendationByBook(ctx context.Context, c *gin.Context)
}

type recommendationRouter struct {
	rh recommendationHandler
}

func New(rh recommendationHandler) *recommendationRouter {
	return &recommendationRouter{rh: rh}
}
