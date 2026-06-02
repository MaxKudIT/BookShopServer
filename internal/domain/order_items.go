package domain

import "github.com/google/uuid"

type OrderItem struct {
	OrderId            uuid.UUID
	PhysicalProductId  uuid.UUID
	Quantity           int
	PriceAtPurchase    float64
	DiscountAtPurchase float64
}
