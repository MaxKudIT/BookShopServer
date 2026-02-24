package fav

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (fs *fStorage) FavByUserId(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	var cartId uuid.UUID
	const GetFavQuery = "SELECT id from favs where user_id = $1"
	if err := fs.db.QueryRowContext(ctx, GetFavQuery, userId).Scan(&cartId); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fs.l.Error("fav not found", "error", err)
			return uuid.Nil, err
		case errors.Is(err, context.Canceled):
			fs.l.Warn("Query cancelled", "error", err)
			return uuid.Nil, err
		case errors.Is(err, context.DeadlineExceeded):
			fs.l.Warn("Query timed out", "error", err)
			return uuid.Nil, err
		default:
			fs.l.Error("Query failed", "error", err)
			return uuid.Nil, err
		}
	}
	fs.l.Info("Successfully got fav", "id", cartId)
	return cartId, nil
}

func (fs *fStorage) Save(ctx context.Context, fav domain.Fav) error {

	const CreateFavQuery = "INSERT INTO favs (id, user_id, created_at) VALUES ($1, $2, $3)"

	if _, err := fs.db.ExecContext(ctx, CreateFavQuery, fav.Id, fav.UserId, fav.CreatedAt); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			fs.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			fs.l.Warn("Query timed out", "error", err)
			return err
		default:
			fs.l.Error("Query failed", "error", err)
			return err
		}
	}
	fs.l.Info("Successfully saved fav", "id", fav.Id)
	return nil
}
