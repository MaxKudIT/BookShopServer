package dto

import "github.com/google/uuid"

type UserSubscriptionDTO struct {
	PlanId uuid.UUID `json:"planId"`
}
