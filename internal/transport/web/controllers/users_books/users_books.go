package users_books

import (
	"context"
	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (ubh *ubHandler) Buy(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var ubdt dto.UsersBooksDTO

	if err := c.ShouldBindJSON(&ubdt); err != nil {
		ubh.l.Error("Error purchase book: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := ubh.ubserv.Buy(ctxnew, ubdt.FirebaseId, ubdt.BookId); err != nil {
		ubh.l.Error("Error purchase book", "bookId", ubdt.BookId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ubh.l.Info("Successfully purchase book", "bookId", ubdt.BookId)
	c.JSON(201, gin.H{"Id:": ubdt.BookId})
}
