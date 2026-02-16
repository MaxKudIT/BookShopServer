package dto

import (
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"time"
)

func CartToDomain(id uuid.UUID, createdAt time.Time) domain.Cart {

	return domain.Cart{
		Id:        id,
		UserId:    uuid.Nil,
		CreatedAt: createdAt,
	}

}
