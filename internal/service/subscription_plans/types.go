package subscription_plans

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type subscriptionPlansStorage interface {
	All(ctx context.Context) ([]domain.SubscriptionPlan, error)
	ById(ctx context.Context, id uuid.UUID) (domain.SubscriptionPlan, error)
	ByTitle(ctx context.Context, title string) (domain.SubscriptionPlan, error)
}

type subscriptionPlansService struct {
	sps subscriptionPlansStorage
	l   *slog.Logger
}

func New(sps subscriptionPlansStorage, l *slog.Logger) *subscriptionPlansService {
	return &subscriptionPlansService{sps: sps, l: l}
}
