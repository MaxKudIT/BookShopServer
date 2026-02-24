package fav

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type favStorage interface {
	Save(ctx context.Context, fav domain.Fav) error
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type favService struct {
	fs favStorage
	us userStorage
	l  *slog.Logger
}

func New(fs favStorage, us userStorage, l *slog.Logger) *favService {
	return &favService{fs: fs, us: us, l: l}
}
