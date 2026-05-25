package book_revs

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type bookRevsStorage interface {
	Save(ctx context.Context, bookReview domain.BookReview) error
	AllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.BookReview, error)
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type bookRevsService struct {
	brs bookRevsStorage
	us  userStorage
	l   *slog.Logger
}

func New(brs bookRevsStorage, us userStorage, l *slog.Logger) *bookRevsService {
	return &bookRevsService{brs: brs, us: us, l: l}
}
