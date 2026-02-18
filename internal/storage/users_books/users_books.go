package users_books

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func (ubs *ubStorage) Buy(ctx context.Context, userId uuid.UUID, bookIds []uuid.UUID) error {

	var CreateUserBookQuery = "INSERT INTO users_books (user_uid, book_id) VALUES "

	placeholders := make([]string, 0, len(bookIds))
	args := make([]any, 0, len(bookIds)*2)

	for i, bookId := range bookIds {
		args = append(args, userId, bookId)

		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
	}

	CreateUserBookQuery += strings.Join(placeholders, ", ")
	ubs.l.Info(CreateUserBookQuery)
	if _, err := ubs.db.ExecContext(ctx, CreateUserBookQuery, args...); err != nil {
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
	ubs.l.Info("Successfully got bought books")
	return nil
}
