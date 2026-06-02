package dto

import (
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"time"
)

type CartItemDTO struct {
	PhysicalBookId uuid.UUID `json:"physicalBookId"`
}

type CartItemsDTO struct {
	PhysicalBookIds []uuid.UUID `json:"physicalBookIds"`
}

type CartItemsDelDTO struct {
	PhysicalBookIds []uuid.UUID `json:"physicalBookIds"`
}

func CartItemToDomain(createdAt time.Time, physicalBookId uuid.UUID) domain.CartItem {
	return domain.CartItem{
		CartId:         uuid.Nil,
		PhysicalBookId: physicalBookId,
		CreatedAt:      createdAt,
	}
}
