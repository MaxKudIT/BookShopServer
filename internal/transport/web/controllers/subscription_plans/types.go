package subscription_plans

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type subscriptionPlansService interface {
	All(ctx context.Context) ([]domain.SubscriptionPlan, error)
	ById(ctx context.Context, id uuid.UUID) (domain.SubscriptionPlan, error)
	ByTitle(ctx context.Context, title string) (domain.SubscriptionPlan, error)
}

type subscriptionPlansHandler struct {
	sps subscriptionPlansService
	l   *slog.Logger
}

func New(sps subscriptionPlansService, l *slog.Logger) *subscriptionPlansHandler {
	return &subscriptionPlansHandler{sps: sps, l: l}
}
