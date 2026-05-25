package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReadingSession struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	BookId    uuid.UUID
	StartedAt time.Time
	EndedAt   *time.Time
	Minutes   int
}
