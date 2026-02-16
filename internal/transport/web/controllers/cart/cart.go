package cart

import (
	"context"
	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (ch *cartHandler) Create(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var cartItemdt struct {
		Title string
	}

	if err := c.ShouldBindJSON(&cartItemdt); err != nil {
		ch.l.Error("Error creating cart: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		ch.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cart := dto.CartToDomain(uuid.New(), time.Now())

	cartId, err := ch.cs.Create(ctxnew, cart, firebaseid.(string))
	if err != nil {
		ch.l.Error("Error creating cart", "id", cartId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ch.l.Info("Successfully created cart", "id", cartId)
	ch.l.Info(cartItemdt.Title)
	c.JSON(201, gin.H{"id": cartId})
}
