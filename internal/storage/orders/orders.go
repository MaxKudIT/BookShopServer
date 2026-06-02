package orders

import (
	"context"
	"errors"

	"github.com/bookshop/internal/domain"
)

func (os *ordersStorage) Save(ctx context.Context, order domain.Order) error {
	const CreateOrderQuery = `
		INSERT INTO orders (id, user_id, status, total_amount, currency, delivery_address, paid_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := os.db.ExecContext(
		ctx,
		CreateOrderQuery,
		order.Id,
		order.UserId,
		order.Status,
		order.TotalAmount,
		order.Currency,
		order.DeliveryAddress,
		order.PaidAt,
	); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			os.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			os.l.Warn("Query timed out", "error", err)
			return err
		default:
			os.l.Error("Query failed", "error", err)
			return err
		}
	}

	os.l.Info("Successfully saved order", "id", order.Id)
	return nil
}
