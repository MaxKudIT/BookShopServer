package domain

import (
	"github.com/google/uuid"
	"time"
)

type FavItem struct {
	FavId     uuid.UUID
	BookId    uuid.UUID
	CreatedAt time.Time
}

type FavItemPreview struct {
	Id       uuid.UUID
	ImageUrl string
	Title    string
	Author   string
	Price    float64
	Discount int
	Rate     float64
}
