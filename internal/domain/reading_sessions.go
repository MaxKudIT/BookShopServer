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

type LastReadingBook struct {
	BookId        uuid.UUID
	LastStartedAt time.Time
}

type ReadingBookPreview struct {
	Id              uuid.UUID
	ImageUrl        string
	Title           string
	Author          string
	Rate            float64
	Genre           Genre
	CreatedDate     time.Time
	PagesCount      int
	ProgressPercent int
	LastStartedAt   time.Time
}
