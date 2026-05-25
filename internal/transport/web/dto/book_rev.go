package dto

import (
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type BookReviewDTO struct {
	BookId uuid.UUID `json:"bookId"`
	Rating float64   `json:"rating"`
}

func BookReviewToDomain(id uuid.UUID, bookReviewDTO BookReviewDTO, createdAt time.Time) domain.BookReview {
	return domain.BookReview{
		Id:        id,
		UserId:    uuid.Nil,
		BookId:    bookReviewDTO.BookId,
		Rating:    bookReviewDTO.Rating,
		CreatedAt: createdAt,
	}
}
