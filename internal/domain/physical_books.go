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

type PhysicalBookStockInfo struct {
	IsInStock     bool
	BookId        uuid.UUID
	Title         string
	Author        string
	Rate          float64
	PhysicalBooks []PhysicalBook
}
