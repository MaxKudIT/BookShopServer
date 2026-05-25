package domain

import "github.com/google/uuid"

type Status string

const (
	Ns       Status = "ns"
	Reading  Status = "reading"
	Finished Status = "finished"
)

type UsersBooks struct {
	UserId          uuid.UUID
	BookId          uuid.UUID
	Status          Status
	ProgressPercent int
	CurrentPage     int
}
