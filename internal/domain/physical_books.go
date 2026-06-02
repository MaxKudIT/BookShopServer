package domain

import "github.com/google/uuid"

type PhysicalBook struct {
	Id         uuid.UUID
	BookId     uuid.UUID
	Price      float64
	Discount   int
	Format     string
	StockCount int
}
