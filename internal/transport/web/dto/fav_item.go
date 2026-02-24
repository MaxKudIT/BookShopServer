package dto

import (
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"time"
)

type FavItemDTO struct {
	BookId uuid.UUID
}

func FavItemToDomain(createdAt time.Time, bookId uuid.UUID) domain.FavItem {
	return domain.FavItem{
		FavId:     uuid.Nil,
		BookId:    bookId,
		CreatedAt: createdAt,
	}
}
