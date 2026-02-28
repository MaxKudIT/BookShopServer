package fav_items

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"strings"
)

func (fis *fiStorage) IsInFavs(ctx context.Context, favId uuid.UUID, bookId uuid.UUID) (bool, error) {
	isInFavs := false
	const IsInFavsQuery = `
    SELECT	
        EXISTS(
        SELECT 1 
        FROM fav_items fi
        WHERE fi.fav_id = $1 
          AND fi.book_id = $2
    )`
	if err := fis.db.QueryRowContext(ctx, IsInFavsQuery, favId, bookId).Scan(
		&isInFavs); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fis.l.Error("book in favs not found", "error", err)
			return false, err
		case errors.Is(err, context.Canceled):
			fis.l.Warn("Query cancelled", "error", err)
			return false, err
		case errors.Is(err, context.DeadlineExceeded):
			fis.l.Warn("Query timed out", "error", err)
			return false, err
		default:
			fis.l.Error("Query failed", "error", err)
			return false, err
		}
	}
	fis.l.Info("Successfully got result about book in the favs", "id", bookId)
	return isInFavs, nil
}

func (fis *fiStorage) Count(ctx context.Context, favId uuid.UUID) (int, error) {

	var count int = 0

	const GetCountFavItemsQuery = "SELECT COUNT(*) FROM fav_items where fav_id = $1"

	if err := fis.db.QueryRowContext(ctx, GetCountFavItemsQuery, favId).Scan(&count); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			fis.l.Warn("Query cancelled", "error", err)
			return 0, err
		case errors.Is(err, context.DeadlineExceeded):
			fis.l.Warn("Query timed out", "error", err)
			return 0, err
		default:
			fis.l.Error("Query failed", "error", err)
			return 0, err
		}
	}
	fis.l.Info("Successfully get count items", "id", favId)
	return count, nil
}

func (fis *fiStorage) AllFavItems(ctx context.Context, favId uuid.UUID) ([]domain.FavItemPreview, error) {

	favItems := make([]domain.FavItemPreview, 0)
	const AllFavItemsQuery = `
    SELECT
        b.id,
        b.title,
        b.author,
        b.price,
        b.discount,
        b.image_url,
        b.rate
    FROM books b
    inner join fav_items fi on fi.book_id = b.id where fi.fav_id = $1
`
	rows, err := fis.db.QueryContext(ctx, AllFavItemsQuery, favId)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			fis.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			fis.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			fis.l.Error("Query failed", "error", err)
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var currentObject domain.FavItemPreview
		if err := rows.Scan(&currentObject.Id, &currentObject.Title, &currentObject.Author, &currentObject.Price, &currentObject.Discount, &currentObject.ImageUrl, &currentObject.Rate); err != nil {
			fis.l.Error("Scan failed", "error", err)
			return nil, err
		}
		favItems = append(favItems, currentObject)
	}

	fis.l.Info("Successfully getting all fav items", "favId", favId)

	return favItems, nil
}

func (fis *fiStorage) Save(ctx context.Context, favItem domain.FavItem) error {

	const CreateFavItemsQuery = "INSERT INTO fav_items (fav_id, book_id, created_at) VALUES ($1, $2, $3)"

	if _, err := fis.db.ExecContext(ctx, CreateFavItemsQuery, favItem.FavId, favItem.BookId, favItem.CreatedAt); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			fis.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			fis.l.Warn("Query timed out", "error", err)
			return err
		default:
			fis.l.Error("Query failed", "error", err)
			return err
		}
	}
	fis.l.Info("Successfully saved fav item", "id", favItem.BookId)
	return nil
}

func (fis *fiStorage) Delete(ctx context.Context, bookIds []uuid.UUID, favId uuid.UUID) error {

	if len(bookIds) == 0 {
		fis.l.Warn("No book IDs provided for deletion", "favId", favId)
		return nil
	}

	var DeleteFavItemQuery = "DELETE FROM fav_items WHERE fav_id = $1 AND book_id IN ("

	args := make([]any, 0, len(bookIds)+1)
	args = append(args, favId)

	placeholders := make([]string, 0, len(bookIds))

	for i := 0; i < len(bookIds); i++ {
		args = append(args, bookIds[i])

		placeholders = append(placeholders, fmt.Sprintf("$%d", i+2))
	}

	DeleteFavItemQuery += strings.Join(placeholders, ", ") + ")"

	if _, err := fis.db.ExecContext(ctx, DeleteFavItemQuery, args...); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fis.l.Error("fav_item not found", "error", err)
			return err
		case errors.Is(err, context.Canceled):
			fis.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			fis.l.Warn("Query timed out", "error", err)
			return err
		default:
			fis.l.Error("Query failed", "error", err)
			return err
		}
	}
	fis.l.Info("Successfully deleted fav_item", "id", favId)
	return nil
}
