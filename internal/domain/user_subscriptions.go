package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionActive        SubscriptionStatus = "active"
	SubscriptionExpired       SubscriptionStatus = "expired"
	SubscriptionCancelled     SubscriptionStatus = "cancelled"
	SubscriptionPaymentFailed SubscriptionStatus = "payment_failed"
)

type UserSubscription struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	PlanId    uuid.UUID
	Status    SubscriptionStatus
	StartedAt time.Time
	ExpiresAt time.Time
}

type UserSubscriptionStatus struct {
	IsActive     bool
	Subscription *UserSubscription
}
