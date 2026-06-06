package physical_books

import (
	"context"

	"github.com/gin-gonic/gin"
)

type physicalBooksHandler interface {
	All(ctx context.Context, c *gin.Context)
	ById(ctx context.Context, c *gin.Context)
	IsPhysicalBookInStock(ctx context.Context, c *gin.Context)
}

type physicalBooksRouter struct {
	pbh physicalBooksHandler
}

func New(pbh physicalBooksHandler) *physicalBooksRouter {
	return &physicalBooksRouter{pbh: pbh}
}
