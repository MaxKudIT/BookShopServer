package domain

import (
	"github.com/google/uuid"
	"time"
)

type Cart struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	CreatedAt time.Time
}
