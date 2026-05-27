package recommendation

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (rh *recommendationHandler) HomeRecommendation(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		rh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit := 10
	if limitQuery := c.Query("limit"); limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)
		if err != nil {
			rh.l.Error("Error parsing recommendation limit", "error", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		limit = parsedLimit
	}

	books, err := rh.rserv.HomeRecommendation(ctxnew, firebaseid.(string), limit)
	if err != nil {
		rh.l.Error("Error getting home recommendations", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rh.l.Info("Successfully got home recommendations")
	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (rh *recommendationHandler) CartRecommendation(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		rh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, ok := rh.recommendationLimit(c)
	if !ok {
		return
	}

	books, err := rh.rserv.CartRecommendation(ctxnew, firebaseid.(string), limit)
	if err != nil {
		rh.l.Error("Error getting cart recommendations", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rh.l.Info("Successfully got cart recommendations")
	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (rh *recommendationHandler) RecommesPage(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		rh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, ok := rh.recommendationLimit(c)
	if !ok {
		return
	}

	forYou, fresh, trend, err := rh.rserv.RecommesPage(ctxnew, firebaseid.(string), limit)
	if err != nil {
		rh.l.Error("Error getting recommendations page", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rh.l.Info("Successfully got recommendations page")
	c.JSON(http.StatusOK, gin.H{
		"forYou": forYou,
		"fresh":  fresh,
		"trend":  trend,
	})
}

func (rh *recommendationHandler) RecommendationByBook(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		rh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bookId, err := uuid.Parse(c.Param("bookId"))
	if err != nil {
		rh.l.Error("Error parsing book id", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, ok := rh.recommendationLimit(c)
	if !ok {
		return
	}

	books, err := rh.rserv.RecommendationByBook(ctxnew, firebaseid.(string), bookId, limit)
	if err != nil {
		rh.l.Error("Error getting recommendations by book", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rh.l.Info("Successfully got recommendations by book", "bookId", bookId)
	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (rh *recommendationHandler) recommendationLimit(c *gin.Context) (int, bool) {
	limit := 10
	if limitQuery := c.Query("limit"); limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)
		if err != nil {
			rh.l.Error("Error parsing recommendation limit", "error", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 0, false
		}
		limit = parsedLimit
	}

	return limit, true
}
