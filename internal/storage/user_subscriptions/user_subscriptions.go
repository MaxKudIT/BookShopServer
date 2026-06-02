package user_subscriptions

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (uss *userSubscriptionsStorage) Save(ctx context.Context, userSubscription domain.UserSubscription) error {
	const CreateUserSubscriptionQuery = `
		INSERT INTO user_subscriptions (id, user_id, plan_id, status, started_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	if _, err := uss.db.ExecContext(
		ctx,
		CreateUserSubscriptionQuery,
		userSubscription.Id,
		userSubscription.UserId,
		userSubscription.PlanId,
		userSubscription.Status,
		userSubscription.StartedAt,
		userSubscription.ExpiresAt,
	); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			uss.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			uss.l.Warn("Query timed out", "error", err)
			return err
		default:
			uss.l.Error("Query failed", "error", err)
			return err
		}
	}

	uss.l.Info("Successfully saved user subscription", "id", userSubscription.Id)
	return nil
}

func (uss *userSubscriptionsStorage) AllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.UserSubscription, error) {
	userSubscriptions := make([]domain.UserSubscription, 0)
	const AllUserSubscriptionsQuery = `
		SELECT id, user_id, plan_id, status, started_at, expires_at
		FROM user_subscriptions
		WHERE user_id = $1
		ORDER BY started_at DESC
	`

	rows, err := uss.db.QueryContext(ctx, AllUserSubscriptionsQuery, userId)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			uss.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			uss.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			uss.l.Error("Query failed", "error", err)
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var currentObject domain.UserSubscription
		if err := rows.Scan(
			&currentObject.Id,
			&currentObject.UserId,
			&currentObject.PlanId,
			&currentObject.Status,
			&currentObject.StartedAt,
			&currentObject.ExpiresAt,
		); err != nil {
			uss.l.Error("Scan failed", "error", err)
			return nil, err
		}
		userSubscriptions = append(userSubscriptions, currentObject)
	}

	if err := rows.Err(); err != nil {
		uss.l.Error("Rows failed", "error", err)
		return nil, err
	}

	uss.l.Info("Successfully got user subscriptions", "userId", userId)
	return userSubscriptions, nil
}

func (uss *userSubscriptionsStorage) ActiveByUserId(ctx context.Context, userId uuid.UUID) (domain.UserSubscription, bool, error) {
	var userSubscription domain.UserSubscription
	const ActiveUserSubscriptionQuery = `
		SELECT id, user_id, plan_id, status, started_at, expires_at
		FROM user_subscriptions
		WHERE user_id = $1
			AND status = 'active'
			AND expires_at > NOW()
		ORDER BY expires_at DESC
		LIMIT 1
	`

	if err := uss.db.QueryRowContext(ctx, ActiveUserSubscriptionQuery, userId).Scan(
		&userSubscription.Id,
		&userSubscription.UserId,
		&userSubscription.PlanId,
		&userSubscription.Status,
		&userSubscription.StartedAt,
		&userSubscription.ExpiresAt,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			uss.l.Info("Active user subscription not found", "userId", userId)
			return domain.UserSubscription{}, false, nil
		case errors.Is(err, context.Canceled):
			uss.l.Warn("Query cancelled", "error", err)
			return domain.UserSubscription{}, false, err
		case errors.Is(err, context.DeadlineExceeded):
			uss.l.Warn("Query timed out", "error", err)
			return domain.UserSubscription{}, false, err
		default:
			uss.l.Error("Query failed", "error", err)
			return domain.UserSubscription{}, false, err
		}
	}

	uss.l.Info("Successfully got active user subscription", "userId", userId)
	return userSubscription, true, nil
}

func (uss *userSubscriptionsStorage) PlanDurationDays(ctx context.Context, planId uuid.UUID) (int, error) {
	var durationDays int
	const PlanDurationDaysQuery = `
		SELECT duration_days
		FROM subscription_plans
		WHERE id = $1
			AND is_active = true
	`

	if err := uss.db.QueryRowContext(ctx, PlanDurationDaysQuery, planId).Scan(&durationDays); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			uss.l.Error("Subscription plan not found", "error", err)
			return 0, err
		case errors.Is(err, context.Canceled):
			uss.l.Warn("Query cancelled", "error", err)
			return 0, err
		case errors.Is(err, context.DeadlineExceeded):
			uss.l.Warn("Query timed out", "error", err)
			return 0, err
		default:
			uss.l.Error("Query failed", "error", err)
			return 0, err
		}
	}

	uss.l.Info("Successfully got subscription plan duration", "planId", planId)
	return durationDays, nil
}
