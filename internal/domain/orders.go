package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderPaid OrderStatus = "paid"
)

type Order struct {
	Id              uuid.UUID
	UserId          uuid.UUID
	Status          OrderStatus
	TotalAmount     float64
	Currency        string
	DeliveryAddress string
	PaidAt          time.Time
}
