package admin

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
)

func (ah *adminHandler) CreateBook(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var bookDTO dto.AdminBookCreateDTO
	if err := c.ShouldBindJSON(&bookDTO); err != nil {
		ah.l.Error("failed to bind admin book payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book payload"})
		return
	}

	result, err := ah.as.CreateBook(ctxnew, firebaseid.(string), domain.AdminBookCreate{
		Title:       bookDTO.Title,
		Author:      bookDTO.Author,
		Genre:       bookDTO.Genre,
		Price:       bookDTO.Price,
		Discount:    bookDTO.Discount,
		ImageUrl:    bookDTO.ImageUrl,
		Description: bookDTO.Description,
		AboutBook:   bookDTO.AboutBook,
		Quote:       bookDTO.Quote,
		ReadingTime: bookDTO.ReadingTime,
		Pages:       bookDTO.Pages,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAdminAccessDenied):
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access denied"})
		case errors.Is(err, domain.ErrInvalidAdminBook):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book payload"})
		default:
			ah.l.Error("failed to create admin book", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		}
		return
	}

	c.JSON(http.StatusCreated, dto.AdminBookCreateResponseDTO{
		BookId:     result.BookId,
		PagesCount: result.PagesCount,
	})
}
