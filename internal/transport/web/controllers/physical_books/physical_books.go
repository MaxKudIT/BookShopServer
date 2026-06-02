package physical_books

import (
	"context"
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
