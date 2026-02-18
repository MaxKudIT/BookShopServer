package cart

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (cs *cStorage) CartByUserId(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	var cartId uuid.UUID
	const GetCartQuery = "SELECT id from carts where user_id = $1"
	if err := cs.db.QueryRowContext(ctx, GetCartQuery, userId).Scan(&cartId); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			cs.l.Error("cart not found", "error", err)
			return uuid.Nil, err
		case errors.Is(err, context.Canceled):
			cs.l.Warn("Query cancelled", "error", err)
			return uuid.Nil, err
		case errors.Is(err, context.DeadlineExceeded):
			cs.l.Warn("Query timed out", "error", err)
			return uuid.Nil, err
		default:
			cs.l.Error("Query failed", "error", err)
			return uuid.Nil, err
		}
	}
	cs.l.Info("Successfully got cart", "id", cartId)
	return cartId, nil
}

func (cs *cStorage) Save(ctx context.Context, cart domain.Cart) error {

	const CreateCartItemsQuery = "INSERT INTO carts (id, user_id, created_at) VALUES ($1, $2, $3)"

	if _, err := cs.db.ExecContext(ctx, CreateCartItemsQuery, cart.Id, cart.UserId, cart.CreatedAt); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			cs.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			cs.l.Warn("Query timed out", "error", err)
			return err
		default:
			cs.l.Error("Query failed", "error", err)
			return err
		}
	}
	cs.l.Info("Successfully saved cart", "id", cart.Id)
	return nil
}
