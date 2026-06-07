package book

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bookshop/internal/domain"
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

func (bh *bookHandler) Search(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		bh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	minPrice, err := optionalFloatQuery(c, "minPrice")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	maxPrice, err := optionalFloatQuery(c, "maxPrice")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	minRate, err := optionalFloatQuery(c, "minRate")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := optionalIntQuery(c, "limit")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offset, err := optionalIntQuery(c, "offset")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := c.Query("q")
	if query == "" {
		query = c.Query("query")
	}

	filter := domain.BookSearchFilter{
		Query:    query,
		Genre:    c.Query("genre"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		MinRate:  minRate,
		Sort:     c.Query("sort"),
		Limit:    limit,
		Offset:   offset,
	}

	books, err := bh.bs.Search(ctxnew, firebaseid.(string), filter)
	if err != nil {
		bh.l.Error("error while searching books", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bh.l.Info("successfully searched books")
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

func optionalFloatQuery(c *gin.Context, key string) (*float64, error) {
	raw := c.Query(key)
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nil, fmt.Errorf("%s must be a number", key)
	}

	return &value, nil
}

func optionalIntQuery(c *gin.Context, key string) (int, error) {
	raw := c.Query(key)
	if raw == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", key)
	}

	return value, nil
}
