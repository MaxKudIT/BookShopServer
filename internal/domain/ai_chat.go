package domain

import (
	"time"

	"github.com/google/uuid"
)

type AIChat struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AIMessage struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	ChatId    uuid.UUID
	Role      string
	Content   string
	CreatedAt time.Time
}
