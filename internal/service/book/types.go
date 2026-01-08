package book

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type bookStorage interface {
	AllBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error)
	AllMyBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error)
	BookById(ctx context.Context, userId uuid.UUID, bookId uuid.UUID) (domain.Book, error)
}

type bookService struct {
	bs bookStorage
	us userStorage
	l  *slog.Logger
}

func New(bs bookStorage, us userStorage, l *slog.Logger) *bookService {
	return &bookService{bs: bs, us: us, l: l}
}
