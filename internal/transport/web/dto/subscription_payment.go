package dto

import (
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionPaymentDTO struct {
	SubId    uuid.UUID `json:"subId"`
	Amount   float64   `json:"amount"`
	Currency string    `json:"currency"`
}

func SubscriptionPaymentToDomain(subscriptionPaymentDTO SubscriptionPaymentDTO) domain.SubscriptionPayment {
	return domain.SubscriptionPayment{
		SubId:    subscriptionPaymentDTO.SubId,
		Amount:   subscriptionPaymentDTO.Amount,
		Currency: subscriptionPaymentDTO.Currency,
	}
}
