package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (us *userStorage) Save(ctx context.Context, user domain.User) error {

	const CREATE_USER_QUERY = "INSERT INTO users (id, login, email, password_hash) VALUES ($1, $2, $3, $4)"

	if _, err := us.db.ExecContext(ctx, CREATE_USER_QUERY, user.Id, user.Login, user.Email, user.PasswordHash); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			us.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			us.l.Warn("Query timed out", "error", err)
			return err
		default:
			us.l.Error("Query failed", "error", err)
			return err
		}
	}
	us.l.Info("Successfully created user", "id", user.Id)
	return nil
}

func (us *userStorage) Delete(ctx context.Context, id uuid.UUID) error {
	const DELETE_USER_QUERY = "DELETE FROM users WHERE id = $1"

	if _, err := us.db.ExecContext(ctx, DELETE_USER_QUERY, id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			us.l.Error("user not found", "error", err)
			return err
		case errors.Is(err, context.Canceled):
			us.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			us.l.Warn("Query timed out", "error", err)
			return err
		default:
			us.l.Error("Query failed", "error", err)
			return err
		}
	}
	us.l.Info("Successfully deleted user", "id", id)
	return nil
}
