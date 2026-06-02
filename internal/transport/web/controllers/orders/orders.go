package orders

import (
	"context"
	"net/http"
	"time"

	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
)

func (oh *ordersHandler) Create(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var orderDTO dto.OrderDTO
	if err := c.ShouldBindJSON(&orderDTO); err != nil {
		oh.l.Error("Error creating order: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		oh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	order := dto.OrderToDomain(orderDTO)
	orderId, err := oh.oserv.Create(ctxnew, firebaseid.(string), order)
	if err != nil {
		oh.l.Error("Error creating order", "id", orderId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oh.l.Info("Successfully created order", "id", orderId)
	c.JSON(http.StatusCreated, gin.H{"id": orderId})
}
