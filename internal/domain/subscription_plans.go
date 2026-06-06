package domain

import "github.com/google/uuid"

type SubscriptionPlan struct {
	Id           uuid.UUID
	Title        string
	Price        float64
	DurationDays int
	IsActive     bool
	Description  string
}
