package book

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (bh *bookHandler) AllBooks(ctx context.Context, c *gin.Context) {
	connew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	//
	//userid, exists := c.Get("user_id")
	//if !exists {
	//	bh.l.Error("userid not found in context")
	//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	//	return
	//}
	//parsingId, err := uuid.Parse(userid.(string))
	parsingId, err := uuid.Parse("a0ef0018-7485-495f-82ee-419948249b16")
	if err != nil {
		//bh.l.Error("error parsing user id", "error", err, "userid", userid)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	books, err := bh.bs.AllBooks(connew, parsingId)
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

	//userid, exists := c.Get("user_id")
	//if !exists {
	//	bh.l.Error("userid not found in context")
	//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	//	return
	//}

	//parsingId, err := uuid.Parse(userid.(string))

	parsingId, err := uuid.Parse("a0ef0018-7485-495f-82ee-419948249b16")
	if err != nil {
		//bh.l.Error("error parsing user id", "error", err, "userid", userid)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	books, err := bh.bs.AllMyBooks(connew, parsingId)
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

	id := c.Param("id")
	parsingId, err := uuid.Parse(id)

	//userid, exists := c.Get("user_id")
	//if !exists {
	//	bh.l.Error("userid not found in context")
	//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	//	return
	//}

	userId, err := uuid.Parse("a0ef0018-7485-495f-82ee-419948249b16")

	if err != nil {
		bh.l.Info("Error parsing id", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	book, err := bh.bs.BookById(ctxnew, userId, parsingId)
	if err != nil {
		bh.l.Info("Error getting book", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	bh.l.Info("Successfully got book", "book", book)
	c.JSON(http.StatusOK, gin.H{"book": book})
}
