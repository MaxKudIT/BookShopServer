package bookviews

import (
	"context"
	"errors"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (bvs *bookViewsStorage) SaveOrUpdate(ctx context.Context, userId uuid.UUID, bookId uuid.UUID) error {
	const CreateBookViewsQuery = `
		INSERT INTO book_views (user_id, book_id, viewed_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id, book_id)
		DO UPDATE SET viewed_at = NOW()
	`

	if _, err := bvs.db.ExecContext(
		ctx,
		CreateBookViewsQuery,
		userId,
		bookId,
	); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			bvs.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			bvs.l.Warn("Query timed out", "error", err)
			return err
		default:
			bvs.l.Error("Query failed", "error", err)
			return err
		}
	}
	bvs.l.Info("Successfully saved book view", "book_id", bookId)
	return nil
}

func (bvs *bookViewsStorage) LastRecords(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookView, error) {
	if limit <= 0 {
		limit = 50
	}

	bookviews := make([]domain.BookView, 0, limit)

	const LastBookViewsByLimit = `
	SELECT user_id, book_id, viewed_at FROM book_views 
	WHERE user_id = $1
	ORDER BY viewed_at DESC
	LIMIT $2
`
	rows, err := bvs.db.QueryContext(ctx, LastBookViewsByLimit, userId, limit)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			bvs.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			bvs.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			bvs.l.Error("Query failed", "error", err)
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var currentObject domain.BookView
		if err := rows.Scan(&currentObject.UserId, &currentObject.BookId, &currentObject.ViewedAt); err != nil {
			bvs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		bookviews = append(bookviews, currentObject)
	}

	if err := rows.Err(); err != nil {
		bvs.l.Error("Rows failed", "error", err)
		return nil, err
	}

	bvs.l.Info("Successfully getting viewed books")

	return bookviews, nil
}
