package page

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

func (ph *pageHandler) AllPagesOfBook(ctx context.Context, c *gin.Context) {
	connew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	bookId := c.Param("id")

	parsingBookId, err := uuid.Parse(bookId)
	if err != nil {
		ph.l.Error("bookId not found in params", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pagesCount, err := ph.ps.AllPagesOfBook(connew, parsingBookId)
	if err != nil {
		ph.l.Error("error while getting count pages", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ph.l.Info("successfully got count pages")
	c.JSON(http.StatusOK, gin.H{"PagesCount": pagesCount})
}

func (ph *pageHandler) PageByNumber(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	bookId := c.Param("id")

	parsingBookId, err := uuid.Parse(bookId)
	if err != nil {
		ph.l.Error("bookId not found in params", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageStr := c.Param("pageNumber")
	if pageStr == "" {
		ph.l.Error("page number not found in params", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageNumber, err := strconv.Atoi(pageStr)
	if err != nil {
		ph.l.Error("page number not formatting in number", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := ph.ps.PageByNumber(ctxnew, pageNumber, parsingBookId)
	if err != nil {
		ph.l.Info("Error getting page", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ph.l.Info("Successfully got page", "page", page.Id)
	c.JSON(http.StatusOK, gin.H{"Page": page})
}
