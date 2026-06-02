package cart_items

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"
)

func (cis *ciStorage) IsInCart(ctx context.Context, cartId uuid.UUID, physicalBookId uuid.UUID) (bool, error) {
	isInCart := false
	const IsInCartQuery = `
    SELECT	
        EXISTS(
        SELECT 1 
        FROM cart_items ci
        WHERE ci.cart_id = $1 
          AND ci.physical_book_id = $2
    )`
	if err := cis.db.QueryRowContext(ctx, IsInCartQuery, cartId, physicalBookId).Scan(
		&isInCart); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			cis.l.Error("book in cart not found", "error", err)
			return false, err
		case errors.Is(err, context.Canceled):
			cis.l.Warn("Query cancelled", "error", err)
			return false, err
		case errors.Is(err, context.DeadlineExceeded):
			cis.l.Warn("Query timed out", "error", err)
			return false, err
		default:
			cis.l.Error("Query failed", "error", err)
			return false, err
		}
	}
	cis.l.Info("Successfully got result about physical book in the cart", "id", physicalBookId)
	return isInCart, nil
}

func (cis *ciStorage) AreAllInCart(ctx context.Context, cartId uuid.UUID, physicalBookIds []uuid.UUID) (bool, error) {
	if len(physicalBookIds) == 0 {
		return false, nil
	}

	IsInCartAllQuery := `
        SELECT COUNT(*) = $1 FROM cart_items
        WHERE cart_id = $2 AND physical_book_id = ANY($3)
    `

	var allExist bool
	if err := cis.db.QueryRowContext(ctx, IsInCartAllQuery, len(physicalBookIds), cartId, pq.Array(physicalBookIds)).Scan(
		&allExist); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			cis.l.Error("books in cart not found", "error", err)
			return false, err
		case errors.Is(err, context.Canceled):
			cis.l.Warn("Query cancelled", "error", err)
			return false, err
		case errors.Is(err, context.DeadlineExceeded):
			cis.l.Warn("Query timed out", "error", err)
			return false, err
		default:
			cis.l.Error("Query failed", "error", err)
			return false, err
		}
	}
	cis.l.Info("Successfully got result about books in the cart")
	return allExist, nil
}

func (cis *ciStorage) AllCartItems(ctx context.Context, cartId uuid.UUID) ([]domain.CartItemPreview, error) {

	cartItems := make([]domain.CartItemPreview, 0)
	const AllCartItemsQuery = `
    SELECT
        pb.id,
		b.id,
        b.title,
        b.author,
        pb.price,
        pb.discount,
        b.image_url,
        b.rate,
		pb.format,
		pb.stock_count
    FROM books b
    INNER JOIN physical_books pb ON pb.book_id = b.id
    INNER JOIN cart_items ci ON ci.physical_book_id = pb.id
	WHERE ci.cart_id = $1
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
		if err := rows.Scan(
			&currentObject.Id,
			&currentObject.BookId,
			&currentObject.Title,
			&currentObject.Author,
			&currentObject.Price,
			&currentObject.Discount,
			&currentObject.ImageUrl,
			&currentObject.Rate,
			&currentObject.Format,
			&currentObject.StockCount,
		); err != nil {
			cis.l.Error("Scan failed", "error", err)
			return nil, err
		}
		cartItems = append(cartItems, currentObject)
	}

	cis.l.Info("Successfully getting all cart items", "cartId", cartId)

	return cartItems, nil
}

func (cis *ciStorage) AllCartItemsId(ctx context.Context, cartId uuid.UUID) ([]uuid.UUID, error) {

	cartItemsId := make([]uuid.UUID, 0)
	const AllCartItemsIdQuery = `
    SELECT physical_book_id from cart_items where cart_id = $1
`
	rows, err := cis.db.QueryContext(ctx, AllCartItemsIdQuery, cartId)
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
		var currentObject uuid.UUID
		if err := rows.Scan(&currentObject); err != nil {
			cis.l.Error("Scan failed", "error", err)
			return nil, err
		}
		cartItemsId = append(cartItemsId, currentObject)
	}

	cis.l.Info("Successfully getting all cart items id", "cartId", cartId)

	return cartItemsId, nil
}

func (cis *ciStorage) Save(ctx context.Context, cartItem domain.CartItem) error {

	const CreateCartItemsQuery = "INSERT INTO cart_items (cart_id, physical_book_id, created_at) VALUES ($1, $2, $3)"

	if _, err := cis.db.ExecContext(ctx, CreateCartItemsQuery, cartItem.CartId, cartItem.PhysicalBookId, cartItem.CreatedAt); err != nil {
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
	cis.l.Info("Successfully saved item", "id", cartItem.PhysicalBookId)
	return nil
}

func (cis *ciStorage) Count(ctx context.Context, cartId uuid.UUID) (int, error) {

	var count int = 0

	const GetCountCartItemsQuery = "SELECT COUNT(*) FROM cart_items where cart_id = $1"

	if err := cis.db.QueryRowContext(ctx, GetCountCartItemsQuery, cartId).Scan(&count); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			cis.l.Warn("Query cancelled", "error", err)
			return 0, err
		case errors.Is(err, context.DeadlineExceeded):
			cis.l.Warn("Query timed out", "error", err)
			return 0, err
		default:
			cis.l.Error("Query failed", "error", err)
			return 0, err
		}
	}
	cis.l.Info("Successfully get count items", "id", cartId)
	return count, nil
}

func (cis *ciStorage) SaveFromFavs(ctx context.Context, cartItems []domain.CartItem) error {
	var CreateSomeCartItemsQuery = "INSERT INTO cart_items (cart_id, physical_book_id, created_at) VALUES "

	placeholders := make([]string, 0, len(cartItems))
	args := make([]any, 0, len(cartItems)*3)

	for i, carti := range cartItems {
		args = append(args, carti.CartId, carti.PhysicalBookId, carti.CreatedAt)

		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
	}

	CreateSomeCartItemsQuery += strings.Join(placeholders, ", ")
	cis.l.Info(CreateSomeCartItemsQuery)
	if _, err := cis.db.ExecContext(ctx, CreateSomeCartItemsQuery, args...); err != nil {
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
	cis.l.Info("Successfully saved items")
	return nil
}

func (cis *ciStorage) Delete(ctx context.Context, physicalBookIds []uuid.UUID, cartId uuid.UUID) error {

	if len(physicalBookIds) == 0 {
		cis.l.Warn("No physical book IDs provided for deletion", "cartId", cartId)
		return nil
	}

	var DeleteCartItemsQuery = "DELETE FROM cart_items WHERE cart_id = $1 AND physical_book_id IN ("

	args := make([]any, 0, len(physicalBookIds)+1)
	args = append(args, cartId)

	placeholders := make([]string, 0, len(physicalBookIds))

	for i := 0; i < len(physicalBookIds); i++ {
		args = append(args, physicalBookIds[i])

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
