package book_revs

import (
	"context"
	"errors"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (brs *bookRevsStorage) Save(ctx context.Context, bookReview domain.BookReview) error {
	const UpdateBookReviewQuery = `
		UPDATE book_reviews
		SET rating = $3,
			created_at = $4
		WHERE user_id = $1
			AND book_id = $2
	`

	result, err := brs.db.ExecContext(
		ctx,
		UpdateBookReviewQuery,
		bookReview.UserId,
		bookReview.BookId,
		bookReview.Rating,
		bookReview.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			brs.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			brs.l.Warn("Query timed out", "error", err)
			return err
		default:
			brs.l.Error("Query failed", "error", err)
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		brs.l.Error("Rows affected failed", "error", err)
		return err
	}
	if rowsAffected > 0 {
		brs.l.Info("Successfully updated book review", "bookId", bookReview.BookId)
		return nil
	}

	const CreateBookReviewQuery = `
		INSERT INTO book_reviews (id, user_id, book_id, rating, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	if _, err := brs.db.ExecContext(
		ctx,
		CreateBookReviewQuery,
		bookReview.Id,
		bookReview.UserId,
		bookReview.BookId,
		bookReview.Rating,
		bookReview.CreatedAt,
	); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			brs.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			brs.l.Warn("Query timed out", "error", err)
			return err
		default:
			brs.l.Error("Query failed", "error", err)
			return err
		}
	}
	brs.l.Info("Successfully saved book review", "id", bookReview.Id)
	return nil
}

func (brs *bookRevsStorage) AllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.BookReview, error) {
	bookReviews := make([]domain.BookReview, 0)
	const AllBookReviewsQuery = `
		SELECT id, user_id, book_id, rating, created_at
		FROM book_reviews
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := brs.db.QueryContext(ctx, AllBookReviewsQuery, userId)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			brs.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			brs.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			brs.l.Error("Query failed", "error", err)
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var currentObject domain.BookReview
		if err := rows.Scan(
			&currentObject.Id,
			&currentObject.UserId,
			&currentObject.BookId,
			&currentObject.Rating,
			&currentObject.CreatedAt,
		); err != nil {
			brs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		bookReviews = append(bookReviews, currentObject)
	}

	if err := rows.Err(); err != nil {
		brs.l.Error("Rows failed", "error", err)
		return nil, err
	}

	brs.l.Info("Successfully got book reviews", "userId", userId)
	return bookReviews, nil
}
