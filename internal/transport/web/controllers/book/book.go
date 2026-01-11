package book

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (bh *bookHandler) AllBooks(ctx context.Context, c *gin.Context) {
	connew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		bh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fmt.Println("firebaseid:", firebaseid)

	books, err := bh.bs.AllBooks(connew, firebaseid.(string))
	if err != nil {
		bh.l.Error("error while getting all books", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bh.l.Info("successfully got all books")
	c.JSON(http.StatusOK, gin.H{"Books": books})
}

func (bh *bookHandler) AllMyBooks(ctx context.Context, c *gin.Context) {
	connew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		bh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fmt.Println("firebaseid:", firebaseid)
	books, err := bh.bs.AllMyBooks(connew, firebaseid.(string))
	if err != nil {
		bh.l.Error("error while getting my books", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bh.l.Info("successfully got my books")
	c.JSON(http.StatusOK, gin.H{"Books": books})
}

func (bh *bookHandler) BookById(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		bh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fmt.Println("firebaseid:", firebaseid)

	id := c.Param("id")
	parsingId, err := uuid.Parse(id)
	if err != nil {
		bh.l.Info("Error parsing book id", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	book, err := bh.bs.BookById(ctxnew, firebaseid.(string), parsingId)
	if err != nil {
		bh.l.Info("Error getting book", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	bh.l.Info("Successfully got book", "book", book)
	c.JSON(http.StatusOK, gin.H{"book": book})
}

func (bh *bookHandler) IsMyBook(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		bh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fmt.Println("firebaseid:", firebaseid)

	id := c.Param("id")
	parsingId, err := uuid.Parse(id)
	if err != nil {
		bh.l.Info("Error parsing book id", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	IsMyBook, err := bh.bs.IsMyBook(ctxnew, firebaseid.(string), parsingId)
	if err != nil {
		bh.l.Info("Error getting result about purchase the book", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	bh.l.Info("Successfully got result", "IsMyBook", IsMyBook)
	c.JSON(http.StatusOK, gin.H{"IsMy": IsMyBook})
}
