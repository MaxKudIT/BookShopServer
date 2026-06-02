package order_items

import (
	"context"
	"net/http"
	"time"

	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (oih *orderItemsHandler) Create(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderId, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		oih.l.Error("Error parsing order id", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var orderItemDTO dto.OrderItemDTO
	if err := c.ShouldBindJSON(&orderItemDTO); err != nil {
		oih.l.Error("Error creating order item: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		oih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orderItem := dto.OrderItemToDomain(orderId, orderItemDTO)
	if err := oih.oiserv.Create(ctxnew, firebaseid.(string), orderItem); err != nil {
		oih.l.Error("Error creating order item", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oih.l.Info("Successfully created order item", "orderId", orderId)
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (oih *orderItemsHandler) All(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderId, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		oih.l.Error("Error parsing order id", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		oih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orderItems, err := oih.oiserv.All(ctxnew, firebaseid.(string), orderId)
	if err != nil {
		oih.l.Error("Error getting order items", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	oih.l.Info("Successfully got order items", "orderId", orderId)
	c.JSON(http.StatusOK, gin.H{"orderItems": orderItems})
}
