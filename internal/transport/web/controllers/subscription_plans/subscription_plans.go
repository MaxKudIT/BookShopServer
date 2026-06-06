package subscription_plans

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (sph *subscriptionPlansHandler) All(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	subscriptionPlans, err := sph.sps.All(ctxnew)
	if err != nil {
		sph.l.Error("Error getting subscription plans", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sph.l.Info("Successfully got subscription plans")
	c.JSON(http.StatusOK, gin.H{"subscriptionPlans": subscriptionPlans})
}

func (sph *subscriptionPlansHandler) ById(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		sph.l.Error("Error parsing subscription plan id", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscriptionPlan, err := sph.sps.ById(ctxnew, id)
	if err != nil {
		sph.l.Error("Error getting subscription plan", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sph.l.Info("Successfully got subscription plan", "id", id)
	c.JSON(http.StatusOK, gin.H{"subscriptionPlan": subscriptionPlan})
}

func (sph *subscriptionPlansHandler) ByTitle(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	title := c.Param("title")
	if title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	subscriptionPlan, err := sph.sps.ByTitle(ctxnew, title)
	if err != nil {
		sph.l.Error("Error getting subscription plan", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sph.l.Info("Successfully got subscription plan", "title", title)
	c.JSON(http.StatusOK, gin.H{"subscriptionPlan": subscriptionPlan})
}
