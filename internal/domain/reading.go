package domain

import "github.com/google/uuid"

type ReadingState struct {
	SessionId       uuid.UUID
	BookId          uuid.UUID
	Status          Status
	CurrentPage     int
	ProgressPercent int
}
