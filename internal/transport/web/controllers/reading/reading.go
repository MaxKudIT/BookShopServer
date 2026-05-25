package reading

import (
	"context"
	"net/http"
	"time"

	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
)

func (rh *readingHandler) Start(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var readingDTO dto.StartReadingDTO
	if err := c.ShouldBindJSON(&readingDTO); err != nil {
		rh.l.Error("Error starting reading: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		rh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	readingState, err := rh.rserv.Start(ctxnew, firebaseid.(string), readingDTO.BookId)
	if err != nil {
		rh.l.Error("Error starting reading", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rh.l.Info("Successfully started reading", "sessionId", readingState.SessionId)
	c.JSON(http.StatusCreated, gin.H{"reading": readingState})
}

func (rh *readingHandler) UpdateProgress(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var readingDTO dto.UpdateReadingProgressDTO
	if err := c.ShouldBindJSON(&readingDTO); err != nil {
		rh.l.Error("Error updating reading progress: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		rh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	readingState, err := rh.rserv.UpdateProgress(ctxnew, firebaseid.(string), readingDTO.BookId, readingDTO.CurrentPage)
	if err != nil {
		rh.l.Error("Error updating reading progress", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rh.l.Info("Successfully updated reading progress", "bookId", readingState.BookId)
	c.JSON(http.StatusOK, gin.H{"reading": readingState})
}

func (rh *readingHandler) Finish(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var readingDTO dto.FinishReadingDTO
	if err := c.ShouldBindJSON(&readingDTO); err != nil {
		rh.l.Error("Error finishing reading: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		rh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	readingState, err := rh.rserv.Finish(ctxnew, firebaseid.(string), readingDTO.SessionId, readingDTO.CurrentPage)
	if err != nil {
		rh.l.Error("Error finishing reading", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rh.l.Info("Successfully finished reading", "sessionId", readingState.SessionId)
	c.JSON(http.StatusOK, gin.H{"reading": readingState})
}
