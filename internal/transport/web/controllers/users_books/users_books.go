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

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		ubh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := ubh.ubserv.Buy(ctxnew, firebaseid.(string), ubdt.BookIds); err != nil {
		ubh.l.Error("Error purchase books", "bookId", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ubh.l.Info("Successfully purchase books")
	c.JSON(201, gin.H{"Ids:": ubdt.BookIds})
	//test
}
