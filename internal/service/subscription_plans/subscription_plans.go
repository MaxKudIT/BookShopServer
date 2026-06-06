package subscription_plans

import (
	"context"
	"strings"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (sps *subscriptionPlansService) All(ctx context.Context) ([]domain.SubscriptionPlan, error) {
	subscriptionPlans, err := sps.sps.All(ctx)
	if err != nil {
		sps.l.Error("Error getting subscription plans", "error", err)
		return nil, err
	}

	sps.l.Info("Successfully got subscription plans")
	return subscriptionPlans, nil
}

func (sps *subscriptionPlansService) ById(ctx context.Context, id uuid.UUID) (domain.SubscriptionPlan, error) {
	subscriptionPlan, err := sps.sps.ById(ctx, id)
	if err != nil {
		sps.l.Error("Error getting subscription plan", "error", err)
		return domain.SubscriptionPlan{}, err
	}

	sps.l.Info("Successfully got subscription plan", "id", id)
	return subscriptionPlan, nil
}

func (sps *subscriptionPlansService) ByTitle(ctx context.Context, title string) (domain.SubscriptionPlan, error) {
	title = strings.TrimSpace(title)
	subscriptionPlan, err := sps.sps.ByTitle(ctx, title)
	if err != nil {
		sps.l.Error("Error getting subscription plan", "error", err)
		return domain.SubscriptionPlan{}, err
	}

	sps.l.Info("Successfully got subscription plan", "title", title)
	return subscriptionPlan, nil
}
