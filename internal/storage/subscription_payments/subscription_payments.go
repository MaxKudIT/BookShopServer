package subscription_payments

import (
	"context"
	"errors"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (sps *subscriptionPaymentsStorage) Save(ctx context.Context, subscriptionPayment domain.SubscriptionPayment) error {
	const CreateSubscriptionPaymentQuery = `
		INSERT INTO subscription_payments (id, user_id, sub_id, amount, currency, paid_at)
		SELECT $1, $2, $3, $4, $5, $6
		WHERE EXISTS (
			SELECT 1
			FROM user_subscriptions
			WHERE id = $3
				AND user_id = $2
		)
	`

	result, err := sps.db.ExecContext(
		ctx,
		CreateSubscriptionPaymentQuery,
		subscriptionPayment.Id,
		subscriptionPayment.UserId,
		subscriptionPayment.SubId,
		subscriptionPayment.Amount,
		subscriptionPayment.Currency,
		subscriptionPayment.PaidAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			sps.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			sps.l.Warn("Query timed out", "error", err)
			return err
		default:
			sps.l.Error("Query failed", "error", err)
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sps.l.Error("Rows affected failed", "error", err)
		return err
	}
	if rowsAffected == 0 {
		err := errors.New("subscription not found for user")
		sps.l.Warn("Subscription payment was not saved", "userId", subscriptionPayment.UserId, "subId", subscriptionPayment.SubId)
		return err
	}

	sps.l.Info("Successfully saved subscription payment", "id", subscriptionPayment.Id)
	return nil
}

func (sps *subscriptionPaymentsStorage) AllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.SubscriptionPayment, error) {
	subscriptionPayments := make([]domain.SubscriptionPayment, 0)
	const AllSubscriptionPaymentsQuery = `
		SELECT id, user_id, sub_id, amount, currency, paid_at
		FROM subscription_payments
		WHERE user_id = $1
		ORDER BY paid_at DESC
	`

	rows, err := sps.db.QueryContext(ctx, AllSubscriptionPaymentsQuery, userId)
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
		var currentObject domain.SubscriptionPayment
		if err := rows.Scan(
			&currentObject.Id,
			&currentObject.UserId,
			&currentObject.SubId,
			&currentObject.Amount,
			&currentObject.Currency,
			&currentObject.PaidAt,
		); err != nil {
			sps.l.Error("Scan failed", "error", err)
			return nil, err
		}
		subscriptionPayments = append(subscriptionPayments, currentObject)
	}

	if err := rows.Err(); err != nil {
		sps.l.Error("Rows failed", "error", err)
		return nil, err
	}

	sps.l.Info("Successfully got subscription payments", "userId", userId)
	return subscriptionPayments, nil
}
