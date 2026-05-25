package stats

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (sh *statsHandler) UserStats(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		sh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userStats, err := sh.sserv.UserStats(ctxnew, firebaseid.(string))
	if err != nil {
		sh.l.Error("Error getting user stats", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sh.l.Info("Successfully got user stats")
	c.JSON(http.StatusOK, gin.H{"stats": userStats})
}
