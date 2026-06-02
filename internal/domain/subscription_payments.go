package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionPayment struct {
	Id       uuid.UUID
	UserId   uuid.UUID
	SubId    uuid.UUID
	Amount   float64
	Currency string
	PaidAt   time.Time
}
