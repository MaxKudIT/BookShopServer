package admin

import (
	"context"
	"database/sql"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (as *adminStorage) IsAdmin(ctx context.Context, firebaseId string) (bool, error) {
	var isAdmin bool

	const query = `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE firebase_id = $1
			  AND role = 'admin'
		)
	`

	if err := as.db.QueryRowContext(ctx, query, firebaseId).Scan(&isAdmin); err != nil {
		as.l.Error("failed to check admin role", "error", err)
		return false, err
	}

	return isAdmin, nil
}

func (as *adminStorage) CreateBookWithPages(ctx context.Context, book domain.AdminBookCreate) (result domain.AdminBookCreateResult, err error) {
	tx, err := as.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		as.l.Error("failed to begin admin book transaction", "error", err)
		return domain.AdminBookCreateResult{}, err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				as.l.Error("failed to rollback admin book transaction", "error", rollbackErr)
			}
		}
	}()

	bookId := uuid.New()
	pagesCount := len(book.Pages)

	const createBookQuery = `
		INSERT INTO books (
			id,
			title,
			pages_count,
			description,
			created_date,
			author,
			genre,
			reading_time,
			price,
			discount,
			about_book,
			quote,
			image_url,
			rate
		) VALUES (
			$1, $2, $3, $4, NOW(), $5, $6, $7, $8, $9, $10, $11, $12, 0
		)
	`

	if _, err = tx.ExecContext(
		ctx,
		createBookQuery,
		bookId,
		book.Title,
		pagesCount,
		book.Description,
		book.Author,
		book.Genre,
		book.ReadingTime,
		book.Price,
		book.Discount,
		book.AboutBook,
		book.Quote,
		book.ImageUrl,
	); err != nil {
		as.l.Error("failed to create admin book", "error", err)
		return domain.AdminBookCreateResult{}, err
	}

	const createPageQuery = `
		INSERT INTO pages (id, number, text, book_id)
		VALUES ($1, $2, $3, $4)
	`

	for index, pageText := range book.Pages {
		if _, err = tx.ExecContext(ctx, createPageQuery, uuid.New(), index+1, pageText, bookId); err != nil {
			as.l.Error("failed to create admin book page", "error", err, "pageNumber", index+1)
			return domain.AdminBookCreateResult{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		as.l.Error("failed to commit admin book transaction", "error", err)
		return domain.AdminBookCreateResult{}, err
	}

	as.l.Info("admin book created", "bookId", bookId, "pagesCount", pagesCount)
	return domain.AdminBookCreateResult{BookId: bookId, PagesCount: pagesCount}, nil
}
