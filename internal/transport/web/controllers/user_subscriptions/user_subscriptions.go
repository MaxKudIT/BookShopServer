package user_subscriptions

import (
	"context"
	"net/http"
	"time"

	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
)

func (ush *userSubscriptionsHandler) Create(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var userSubscriptionDTO dto.UserSubscriptionDTO
	if err := c.ShouldBindJSON(&userSubscriptionDTO); err != nil {
		ush.l.Error("Error creating user subscription: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		ush.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userSubscriptionId, err := ush.usserv.Create(ctxnew, firebaseid.(string), userSubscriptionDTO.PlanId)
	if err != nil {
		ush.l.Error("Error creating user subscription", "id", userSubscriptionId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ush.l.Info("Successfully created user subscription", "id", userSubscriptionId)
	c.JSON(http.StatusCreated, gin.H{"id": userSubscriptionId})
}

func (ush *userSubscriptionsHandler) All(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		ush.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userSubscriptions, err := ush.usserv.All(ctxnew, firebaseid.(string))
	if err != nil {
		ush.l.Error("Error getting user subscriptions", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ush.l.Info("Successfully got user subscriptions")
	c.JSON(http.StatusOK, gin.H{"subscriptions": userSubscriptions})
}

func (ush *userSubscriptionsHandler) Status(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		ush.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	status, err := ush.usserv.Status(ctxnew, firebaseid.(string))
	if err != nil {
		ush.l.Error("Error getting user subscription status", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ush.l.Info("Successfully got user subscription status")
	c.JSON(http.StatusOK, gin.H{"status": status})
}
