package subscription_payments

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type subscriptionPaymentsService interface {
	Create(ctx context.Context, firebaseId string, subscriptionPayment domain.SubscriptionPayment) (uuid.UUID, error)
	All(ctx context.Context, firebaseId string) ([]domain.SubscriptionPayment, error)
}

type subscriptionPaymentsHandler struct {
	spserv subscriptionPaymentsService
	l      *slog.Logger
}

func New(spserv subscriptionPaymentsService, l *slog.Logger) *subscriptionPaymentsHandler {
	return &subscriptionPaymentsHandler{spserv: spserv, l: l}
}
