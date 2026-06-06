package physical_books

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (pbh *physicalBooksHandler) All(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	physicalBooks, err := pbh.pbserv.All(ctxnew)
	if err != nil {
		pbh.l.Error("Error getting physical books", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pbh.l.Info("Successfully got physical books")
	c.JSON(http.StatusOK, gin.H{"physicalBooks": physicalBooks})
}

func (pbh *physicalBooksHandler) ById(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pbh.l.Error("Error parsing physical book id", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	physicalBook, err := pbh.pbserv.ById(ctxnew, id)
	if err != nil {
		pbh.l.Error("Error getting physical book", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pbh.l.Info("Successfully got physical book", "id", id)
	c.JSON(http.StatusOK, gin.H{"physicalBook": physicalBook})
}

func (pbh *physicalBooksHandler) IsPhysicalBookInStock(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	bookId, err := uuid.Parse(c.Param("bookId"))
	if err != nil {
		pbh.l.Error("Error parsing book id", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	physicalBookStockInfo, err := pbh.pbserv.IsPhysicalBookInStock(ctxnew, bookId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			pbh.l.Error("Book not found", "error", err)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}

		pbh.l.Error("Error getting physical book stock info", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pbh.l.Info("Successfully got physical book stock info", "bookId", bookId)
	c.JSON(http.StatusOK, gin.H{"physicalBookStockInfo": physicalBookStockInfo})
}
