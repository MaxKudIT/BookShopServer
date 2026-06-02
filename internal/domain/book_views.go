package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookView struct {
	UserId   uuid.UUID
	BookId   uuid.UUID
	ViewedAt time.Time
}

type BookViewPreview struct {
	Id          uuid.UUID
	ImageUrl    string
	Title       string
	Author      string
	CreatedDate time.Time
	Genre       Genre
	Rate        float64
	PagesCount  int
	ViewedAt    time.Time
}
