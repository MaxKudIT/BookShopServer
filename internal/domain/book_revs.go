package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookReview struct {
	Id        uuid.UUID
	UserId    uuid.UUID // составной
	BookId    uuid.UUID // ключ
	Rating    float64
	CreatedAt time.Time
}
