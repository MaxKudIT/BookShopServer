package domain

import (
	"github.com/google/uuid"
	"time"
)

type CartItem struct {
	CartId    uuid.UUID
	BookId    uuid.UUID
	CreatedAt time.Time
}

type CartItemPreview struct {
	Id       uuid.UUID
	ImageUrl string
	Title    string
	Author   string
	Price    float64
	Discount int
	Rate     float64
}
