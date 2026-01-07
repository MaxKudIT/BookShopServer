package domain

import "github.com/google/uuid"

type Page struct {
	Id     uuid.UUID
	Number int
	Text   string
}
