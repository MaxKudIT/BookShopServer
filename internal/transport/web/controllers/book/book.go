package book

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	//select {
	//case <-connew.Done():
	//	if errors.Is(connew.Err(), context.DeadlineExceeded) { 504
	//		c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
	//			"error": "Превышено время ожидания",
	//		})
	//		return
	//	}
	//default:
	//
	//}

	//if true { //иммитация падения сервиса
	//	c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{ 503
	//		"error":       "Сервис временно недоступен",
	//		"retry_after": 30,
	//	})
	//	return
	//}

	books, err := bh.bs.AllMyBooks(connew, firebaseid.(string))
	if err != nil {
		bh.l.Error("error while getting my books", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //500
		return
	}

	//_, err = http.Get("https://external-api.com") 502
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
	//		"error": "Ошибка при обращении к внешнему сервису",
	//	})
	//	return
	//}

	//c.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{ 501
	//	"error":   "Метод не реализован",
	//	"method":  c.Request.Method,
	//	"path":    c.Request.URL.Path,
	//	"message": "Этот функционал появится в следующем обновлении",
	//})

	bh.l.Info("successfully got my books")
	c.JSON(http.StatusOK, gin.H{"Books": books})
}

func (bh *bookHandler) AllNotMyBooks(ctx context.Context, c *gin.Context) {
	connew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		bh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fmt.Println("firebaseid:", firebaseid)

	books, err := bh.bs.AllNotMyBooks(connew, firebaseid.(string))
	if err != nil {
		bh.l.Error("error while getting not my books", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //500
		return
	}

	bh.l.Info("successfully got not my books")
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
