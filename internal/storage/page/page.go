package page

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (ps *pageStorage) AllPagesOfBook(ctx context.Context, bookId uuid.UUID) (int, error) {

	count := 0

	const AllPagesQuery = `
    SELECT COUNT(*)
    FROM pages 
    where book_id = $1
`
	if err := ps.db.QueryRowContext(ctx, AllPagesQuery, bookId).Scan(&count); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			ps.l.Warn("Query cancelled", "error", err)
			return 0, err
		case errors.Is(err, context.DeadlineExceeded):
			ps.l.Warn("Query timed out", "error", err)
			return 0, err
		default:
			ps.l.Error("Query failed", "error", err)
			return 0, err
		}
	}

	ps.l.Info("Successfully getting count pages", "bookid", bookId)

	return count, nil
}

func (ps *pageStorage) PageByNumber(ctx context.Context, pageNumber int, bookId uuid.UUID) (domain.Page, error) {
	page := domain.Page{}
	const GetPagesQuery = "SELECT id, number, text from pages where number = $1 AND book_id = $2"
	if err := ps.db.QueryRowContext(ctx, GetPagesQuery, pageNumber, bookId).Scan(
		&page.Id,
		&page.Number,
		&page.Text); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ps.l.Error("page not found", "error", err)
			return domain.Page{}, err
		case errors.Is(err, context.Canceled):
			ps.l.Warn("Query cancelled", "error", err)
			return domain.Page{}, err
		case errors.Is(err, context.DeadlineExceeded):
			ps.l.Warn("Query timed out", "error", err)
			return domain.Page{}, err
		default:
			ps.l.Error("Query failed", "error", err)
			return domain.Page{}, err
		}
	}
	ps.l.Info("Successfully got page", "id", page.Id)
	return page, nil
}
