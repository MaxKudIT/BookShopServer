package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (us *userStorage) UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error) {

	var userId uuid.UUID

	const GetUserIdQuery = "SELECT id from users where firebase_id = $1"
	if err := us.db.QueryRowContext(ctx, GetUserIdQuery, firebaseId).Scan(&userId); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			us.l.Error("user not found", "error", err)
			return uuid.Nil, err
		case errors.Is(err, context.Canceled):
			us.l.Warn("Query cancelled", "error", err)
			return uuid.Nil, err
		case errors.Is(err, context.DeadlineExceeded):
			us.l.Warn("Query timed out", "error", err)
			return uuid.Nil, err
		default:
			us.l.Error("Query failed", "error", err)
			return uuid.Nil, err
		}
	}
	us.l.Info("Successfully got id", "id")
	return userId, nil
}

func (us *userStorage) Save(ctx context.Context, user domain.User) error {

	const CreateUserQuery = "INSERT INTO users (id, firebase_id) VALUES ($1, $2)"

	if _, err := us.db.ExecContext(ctx, CreateUserQuery, user.Id, user.FirebaseId); err != nil {
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
	const DeleteUserQuery = "DELETE FROM users WHERE id = $1"

	if _, err := us.db.ExecContext(ctx, DeleteUserQuery, id); err != nil {
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
