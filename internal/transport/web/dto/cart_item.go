package dto

import (
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"time"
)

type CartItemDTO struct {
	BookId uuid.UUID
}

type CartItemsDTO struct {
	BookIds []uuid.UUID
}

type CartItemsDelDTO struct {
	BookIds []uuid.UUID
}

func CartItemToDomain(createdAt time.Time, bookId uuid.UUID) domain.CartItem {
	return domain.CartItem{
		CartId:    uuid.Nil,
		BookId:    bookId,
		CreatedAt: createdAt,
	}
}
