package users_books

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

func (ubs *ubStorage) Buy(ctx context.Context, userId uuid.UUID, bookId uuid.UUID) error {

	const CreateUserBookQuery = "INSERT INTO users_books (user_uid, book_id) VALUES ($1, $2)"

	if _, err := ubs.db.ExecContext(ctx, CreateUserBookQuery, userId, bookId); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			ubs.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			ubs.l.Warn("Query timed out", "error", err)
			return err
		default:
			ubs.l.Error("Query failed", "error", err)
			return err
		}
	}
	ubs.l.Info("Successfully got bought book", "idBook", bookId)
	return nil
}
