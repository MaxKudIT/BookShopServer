package subscription_plans

import (
	"context"
	"errors"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (sps *subscriptionPlansStorage) All(ctx context.Context) ([]domain.SubscriptionPlan, error) {
	subscriptionPlans := make([]domain.SubscriptionPlan, 0)
	const AllSubscriptionPlansQuery = `
		SELECT id, title, price, duration_days, is_active, description
		FROM subscription_plans
		WHERE is_active = true
		ORDER BY duration_days ASC
	`

	rows, err := sps.db.QueryContext(ctx, AllSubscriptionPlansQuery)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			sps.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			sps.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			sps.l.Error("Query failed", "error", err)
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var subscriptionPlan domain.SubscriptionPlan
		if err := rows.Scan(
			&subscriptionPlan.Id,
			&subscriptionPlan.Title,
			&subscriptionPlan.Price,
			&subscriptionPlan.DurationDays,
			&subscriptionPlan.IsActive,
			&subscriptionPlan.Description,
		); err != nil {
			sps.l.Error("Scan failed", "error", err)
			return nil, err
		}
		subscriptionPlans = append(subscriptionPlans, subscriptionPlan)
	}

	if err := rows.Err(); err != nil {
		sps.l.Error("Rows failed", "error", err)
		return nil, err
	}

	sps.l.Info("Successfully got subscription plans")
	return subscriptionPlans, nil
}

func (sps *subscriptionPlansStorage) ById(ctx context.Context, id uuid.UUID) (domain.SubscriptionPlan, error) {
	var subscriptionPlan domain.SubscriptionPlan
	const SubscriptionPlanByIdQuery = `
		SELECT id, title, price, duration_days, is_active, description
		FROM subscription_plans
		WHERE id = $1
			AND is_active = true
	`

	if err := sps.db.QueryRowContext(ctx, SubscriptionPlanByIdQuery, id).Scan(
		&subscriptionPlan.Id,
		&subscriptionPlan.Title,
		&subscriptionPlan.Price,
		&subscriptionPlan.DurationDays,
		&subscriptionPlan.IsActive,
		&subscriptionPlan.Description,
	); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			sps.l.Warn("Query cancelled", "error", err)
			return domain.SubscriptionPlan{}, err
		case errors.Is(err, context.DeadlineExceeded):
			sps.l.Warn("Query timed out", "error", err)
			return domain.SubscriptionPlan{}, err
		default:
			sps.l.Error("Query failed", "error", err)
			return domain.SubscriptionPlan{}, err
		}
	}

	sps.l.Info("Successfully got subscription plan", "id", id)
	return subscriptionPlan, nil
}

func (sps *subscriptionPlansStorage) ByTitle(ctx context.Context, title string) (domain.SubscriptionPlan, error) {
	var subscriptionPlan domain.SubscriptionPlan
	const SubscriptionPlanByTitleQuery = `
		SELECT id, title, price, duration_days, is_active, description
		FROM subscription_plans
		WHERE title = $1
			AND is_active = true
	`

	if err := sps.db.QueryRowContext(ctx, SubscriptionPlanByTitleQuery, title).Scan(
		&subscriptionPlan.Id,
		&subscriptionPlan.Title,
		&subscriptionPlan.Price,
		&subscriptionPlan.DurationDays,
		&subscriptionPlan.IsActive,
		&subscriptionPlan.Description,
	); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			sps.l.Warn("Query cancelled", "error", err)
			return domain.SubscriptionPlan{}, err
		case errors.Is(err, context.DeadlineExceeded):
			sps.l.Warn("Query timed out", "error", err)
			return domain.SubscriptionPlan{}, err
		default:
			sps.l.Error("Query failed", "error", err)
			return domain.SubscriptionPlan{}, err
		}
	}

	sps.l.Info("Successfully got subscription plan", "title", title)
	return subscriptionPlan, nil
}
