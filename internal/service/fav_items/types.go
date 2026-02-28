package fav_items

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type favStorage interface {
	FavByUserId(ctx context.Context, userId uuid.UUID) (uuid.UUID, error)
}

type favItemsStorage interface {
	IsInFavs(ctx context.Context, favId uuid.UUID, bookId uuid.UUID) (bool, error)
	AllFavItems(ctx context.Context, favId uuid.UUID) ([]domain.FavItemPreview, error)
	Save(ctx context.Context, favItem domain.FavItem) error
	Count(ctx context.Context, favId uuid.UUID) (int, error)
	Delete(ctx context.Context, bookIds []uuid.UUID, favId uuid.UUID) error
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type favItemsService struct {
	fs  favStorage
	fis favItemsStorage
	us  userStorage
	l   *slog.Logger
}

func New(fis favItemsStorage, fs favStorage, us userStorage, l *slog.Logger) *favItemsService {
	return &favItemsService{fis: fis, fs: fs, us: us, l: l}
}
