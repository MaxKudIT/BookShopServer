package dto

import "github.com/google/uuid"

type StartReadingDTO struct {
	BookId uuid.UUID `json:"bookId"`
}

type UpdateReadingProgressDTO struct {
	BookId      uuid.UUID `json:"bookId"`
	CurrentPage int       `json:"currentPage"`
}

type FinishReadingDTO struct {
	SessionId   uuid.UUID `json:"sessionId"`
	CurrentPage int       `json:"currentPage"`
}
