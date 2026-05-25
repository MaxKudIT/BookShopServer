package book_revs

import (
	"context"
	"net/http"
	"time"

	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (brh *bookRevsHandler) Create(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var bookReviewDTO dto.BookReviewDTO
	if err := c.ShouldBindJSON(&bookReviewDTO); err != nil {
		brh.l.Error("Error creating book review: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		brh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bookReview := dto.BookReviewToDomain(uuid.New(), bookReviewDTO, time.Now())
	bookReviewId, err := brh.brserv.Create(ctxnew, bookReview, firebaseid.(string))
	if err != nil {
		brh.l.Error("Error creating book review", "id", bookReviewId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	brh.l.Info("Successfully created book review", "id", bookReviewId)
	c.JSON(http.StatusCreated, gin.H{"id": bookReviewId})
}

func (brh *bookRevsHandler) All(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		brh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bookReviews, err := brh.brserv.All(ctxnew, firebaseid.(string))
	if err != nil {
		brh.l.Error("Error getting book reviews", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	brh.l.Info("Successfully got book reviews")
	c.JSON(http.StatusOK, gin.H{"bookReviews": bookReviews})
}
