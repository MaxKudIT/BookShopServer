package reading_sessions

import (
	"context"
	"net/http"
	"time"

	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (rsh *readingSessionsHandler) Create(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var readingSessionDTO dto.ReadingSessionDTO
	if err := c.ShouldBindJSON(&readingSessionDTO); err != nil {
		rsh.l.Error("Error creating reading session: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		rsh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	readingSession := dto.ReadingSessionToDomain(uuid.New(), readingSessionDTO)
	readingSessionId, err := rsh.rsserv.Create(ctxnew, readingSession, firebaseid.(string))
	if err != nil {
		rsh.l.Error("Error creating reading session", "id", readingSessionId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rsh.l.Info("Successfully created reading session", "id", readingSessionId)
	c.JSON(http.StatusCreated, gin.H{"id": readingSessionId})
}

func (rsh *readingSessionsHandler) All(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		rsh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	readingSessions, err := rsh.rsserv.All(ctxnew, firebaseid.(string))
	if err != nil {
		rsh.l.Error("Error getting reading sessions", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsh.l.Info("Successfully got reading sessions")
	c.JSON(http.StatusOK, gin.H{"readingSessions": readingSessions})
}
