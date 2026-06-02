package subscription_payments

import (
	"context"
	"net/http"
	"time"

	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
)

func (sph *subscriptionPaymentsHandler) Create(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var subscriptionPaymentDTO dto.SubscriptionPaymentDTO
	if err := c.ShouldBindJSON(&subscriptionPaymentDTO); err != nil {
		sph.l.Error("Error creating subscription payment: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		sph.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subscriptionPayment := dto.SubscriptionPaymentToDomain(subscriptionPaymentDTO)
	subscriptionPaymentId, err := sph.spserv.Create(ctxnew, firebaseid.(string), subscriptionPayment)
	if err != nil {
		sph.l.Error("Error creating subscription payment", "id", subscriptionPaymentId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sph.l.Info("Successfully created subscription payment", "id", subscriptionPaymentId)
	c.JSON(http.StatusCreated, gin.H{"id": subscriptionPaymentId})
}

func (sph *subscriptionPaymentsHandler) All(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		sph.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subscriptionPayments, err := sph.spserv.All(ctxnew, firebaseid.(string))
	if err != nil {
		sph.l.Error("Error getting subscription payments", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sph.l.Info("Successfully got subscription payments")
	c.JSON(http.StatusOK, gin.H{"payments": subscriptionPayments})
}
