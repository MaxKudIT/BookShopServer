package order_items

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (ois *orderItemsStorage) Save(ctx context.Context, userId uuid.UUID, orderItem domain.OrderItem) error {
	tx, err := ois.db.BeginTx(ctx, nil)
	if err != nil {
		ois.l.Error("Begin transaction failed", "error", err)
		return err
	}
	defer tx.Rollback()

	var orderExists bool
	const CheckOrderQuery = `
		SELECT EXISTS(
			SELECT 1
			FROM orders
			WHERE id = $1
				AND user_id = $2
		)
	`
	if err := tx.QueryRowContext(ctx, CheckOrderQuery, orderItem.OrderId, userId).Scan(&orderExists); err != nil {
		ois.l.Error("Order check failed", "error", err)
		return err
	}
	if !orderExists {
		err := errors.New("order not found for user")
		ois.l.Warn("Order not found for user", "userId", userId, "orderId", orderItem.OrderId)
		return err
	}

	const ReservePhysicalBookQuery = `
		UPDATE physical_books
		SET stock_count = stock_count - $2
		WHERE id = $1
			AND stock_count >= $2
		RETURNING price, discount
	`
	if err := tx.QueryRowContext(ctx, ReservePhysicalBookQuery, orderItem.PhysicalProductId, orderItem.Quantity).Scan(
		&orderItem.PriceAtPurchase,
		&orderItem.DiscountAtPurchase,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ois.l.Warn("Physical book not found or not enough stock", "physicalProductId", orderItem.PhysicalProductId)
			return err
		default:
			ois.l.Error("Physical book reservation failed", "error", err)
			return err
		}
	}

	const UpdateOrderItemQuery = `
		UPDATE order_items
		SET quantity = quantity + $3,
			price_at_purchase = $4,
			discount_at_purchase = $5
		WHERE order_id = $1
			AND ph_product_id = $2
	`
	result, err := tx.ExecContext(
		ctx,
		UpdateOrderItemQuery,
		orderItem.OrderId,
		orderItem.PhysicalProductId,
		orderItem.Quantity,
		orderItem.PriceAtPurchase,
		orderItem.DiscountAtPurchase,
	)
	if err != nil {
		ois.l.Error("Order item update failed", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ois.l.Error("Rows affected failed", "error", err)
		return err
	}

	if rowsAffected == 0 {
		const CreateOrderItemQuery = `
			INSERT INTO order_items (
				order_id,
				ph_product_id,
				quantity,
				price_at_purchase,
				discount_at_purchase
			)
			VALUES ($1, $2, $3, $4, $5)
		`
		if _, err := tx.ExecContext(
			ctx,
			CreateOrderItemQuery,
			orderItem.OrderId,
			orderItem.PhysicalProductId,
			orderItem.Quantity,
			orderItem.PriceAtPurchase,
			orderItem.DiscountAtPurchase,
		); err != nil {
			ois.l.Error("Order item insert failed", "error", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		ois.l.Error("Commit transaction failed", "error", err)
		return err
	}

	ois.l.Info("Successfully saved order item", "orderId", orderItem.OrderId, "physicalProductId", orderItem.PhysicalProductId)
	return nil
}

func (ois *orderItemsStorage) AllByOrderId(ctx context.Context, userId uuid.UUID, orderId uuid.UUID) ([]domain.OrderItem, error) {
	orderItems := make([]domain.OrderItem, 0)
	const AllOrderItemsQuery = `
		SELECT
			oi.order_id,
			oi.ph_product_id,
			oi.quantity,
			oi.price_at_purchase,
			oi.discount_at_purchase
		FROM order_items oi
		INNER JOIN orders o ON o.id = oi.order_id
		WHERE oi.order_id = $1
			AND o.user_id = $2
		ORDER BY oi.physical_product_id
	`

	rows, err := ois.db.QueryContext(ctx, AllOrderItemsQuery, orderId, userId)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			ois.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			ois.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			ois.l.Error("Query failed", "error", err)
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var currentObject domain.OrderItem
		if err := rows.Scan(
			&currentObject.OrderId,
			&currentObject.PhysicalProductId,
			&currentObject.Quantity,
			&currentObject.PriceAtPurchase,
			&currentObject.DiscountAtPurchase,
		); err != nil {
			ois.l.Error("Scan failed", "error", err)
			return nil, err
		}
		orderItems = append(orderItems, currentObject)
	}

	if err := rows.Err(); err != nil {
		ois.l.Error("Rows failed", "error", err)
		return nil, err
	}

	ois.l.Info("Successfully got order items", "orderId", orderId)
	return orderItems, nil
}
