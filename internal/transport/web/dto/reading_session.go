package dto

import (
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type ReadingSessionDTO struct {
	BookId    uuid.UUID  `json:"bookId"`
	StartedAt time.Time  `json:"startedAt"`
	EndedAt   *time.Time `json:"endedAt"`
	Minutes   int        `json:"minutes"`
}

func ReadingSessionToDomain(id uuid.UUID, readingSessionDTO ReadingSessionDTO) domain.ReadingSession {
	return domain.ReadingSession{
		Id:        id,
		UserId:    uuid.Nil,
		BookId:    readingSessionDTO.BookId,
		StartedAt: readingSessionDTO.StartedAt,
		EndedAt:   readingSessionDTO.EndedAt,
		Minutes:   readingSessionDTO.Minutes,
	}
}
