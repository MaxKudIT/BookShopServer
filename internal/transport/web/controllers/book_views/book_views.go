package bookviews

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (bvh *bookViewsHandler) SaveOrUpdate(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		bvh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bookId, err := uuid.Parse(c.Param("bookId"))
	if err != nil {
		bvh.l.Error("Error parsing book id", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bvh.bvserv.SaveOrUpdate(ctxnew, firebaseid.(string), bookId); err != nil {
		bvh.l.Error("Error saving book view", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bvh.l.Info("Successfully saved book view", "bookId", bookId)
	c.JSON(http.StatusCreated, gin.H{"bookId": bookId})
}

func (bvh *bookViewsHandler) LastRecords(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		bvh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit := 50
	if limitQuery := c.Query("limit"); limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)
		if err != nil {
			bvh.l.Error("Error parsing book views limit", "error", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		limit = parsedLimit
	}

	bookPreviews, err := bvh.bvserv.LastRecords(ctxnew, firebaseid.(string), limit)
	if err != nil {
		bvh.l.Error("Error getting book views", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bvh.l.Info("Successfully got book views")
	c.JSON(http.StatusOK, gin.H{"books": bookPreviews})
}
