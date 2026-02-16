package cart_items

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"strings"
)

func (cis *ciStorage) AllCartItems(ctx context.Context, cartId uuid.UUID) ([]domain.CartItemPreview, error) {

	cartItems := make([]domain.CartItemPreview, 0)
	const AllCartItemsQuery = `
    SELECT
        b.id,
        b.title,
        b.author,
        b.price,
        b.discount,
        b.image_url,
        b.rate
    FROM books b
    inner join cart_items ci on ci.book_id = b.id where ci.cart_id = $1
`
	rows, err := cis.db.QueryContext(ctx, AllCartItemsQuery, cartId)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			cis.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			cis.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			cis.l.Error("Query failed", "error", err)
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var currentObject domain.CartItemPreview
		if err := rows.Scan(&currentObject.Id, &currentObject.Title, &currentObject.Author, &currentObject.Price, &currentObject.Discount, &currentObject.ImageUrl, &currentObject.Rate); err != nil {
			cis.l.Error("Scan failed", "error", err)
			return nil, err
		}
		cartItems = append(cartItems, currentObject)
	}

	cis.l.Info("Successfully getting all cart items", "cartId", cartId)

	return cartItems, nil
}

func (cis *ciStorage) Save(ctx context.Context, cartItem domain.CartItem) error {

	const CreateCartItemsQuery = "INSERT INTO cart_items (cart_id, book_id, created_at) VALUES ($1, $2, $3)"

	if _, err := cis.db.ExecContext(ctx, CreateCartItemsQuery, cartItem.CartId, cartItem.BookId, cartItem.CreatedAt); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			cis.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			cis.l.Warn("Query timed out", "error", err)
			return err
		default:
			cis.l.Error("Query failed", "error", err)
			return err
		}
	}
	cis.l.Info("Successfully saved item", "id", cartItem.BookId)
	return nil
}

func (cis *ciStorage) Delete(ctx context.Context, bookIds []uuid.UUID, cartId uuid.UUID) error {

	var DeleteCartItemsQuery = "DELETE FROM cart_items WHERE cart_id = $1 AND book_id IN ("

	args := make([]any, 0, len(bookIds)+1)
	args = append(args, cartId)

	placeholders := make([]string, 0, len(bookIds))

	for i := 0; i < len(bookIds); i++ {
		args = append(args, bookIds[i])

		placeholders = append(placeholders, fmt.Sprintf("$%d", i+2))
	}

	DeleteCartItemsQuery += strings.Join(placeholders, ", ") + ")"

	if _, err := cis.db.ExecContext(ctx, DeleteCartItemsQuery, args...); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			cis.l.Error("cart_item not found", "error", err)
			return err
		case errors.Is(err, context.Canceled):
			cis.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			cis.l.Warn("Query timed out", "error", err)
			return err
		default:
			cis.l.Error("Query failed", "error", err)
			return err
		}
	}
	cis.l.Info("Successfully deleted cart_item", "id", cartId)
	return nil
}
