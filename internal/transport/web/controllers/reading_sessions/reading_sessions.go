package reading_sessions

import (
	"context"
	"net/http"
	"strconv"
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

func (rsh *readingSessionsHandler) Close(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	parsingSessionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		rsh.l.Error("Error parsing reading session id", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		rsh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	readingSession, err := rsh.rsserv.Close(ctxnew, firebaseid.(string), parsingSessionId)
	if err != nil {
		rsh.l.Error("Error closing reading session", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rsh.l.Info("Successfully closed reading session", "id", parsingSessionId)
	c.JSON(http.StatusOK, gin.H{"readingSession": readingSession})
}

func (rsh *readingSessionsHandler) LastReadingBooks(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		rsh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit := 4
	if limitQuery := c.Query("limit"); limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)
		if err != nil {
			rsh.l.Error("Error parsing last reading books limit", "error", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		limit = parsedLimit
	}

	books, err := rsh.rsserv.LastReadingBooks(ctxnew, firebaseid.(string), limit)
	if err != nil {
		rsh.l.Error("Error getting last reading books", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsh.l.Info("Successfully got last reading books")
	c.JSON(http.StatusOK, gin.H{"books": books})
}
