package dto

import (
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"time"
)

func FavToDomain(id uuid.UUID, createdAt time.Time) domain.Fav {

	return domain.Fav{
		Id:        id,
		UserId:    uuid.Nil,
		CreatedAt: createdAt,
	}

}
