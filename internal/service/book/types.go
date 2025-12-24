package book

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type bookStorage interface {
	AllBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error)
	AllMyBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error)
	BookById(ctx context.Context, id uuid.UUID) (domain.Book, error)
}

type bookService struct {
	bs bookStorage
	l  *slog.Logger
}

func New(bs bookStorage, l *slog.Logger) *bookService {
	return &bookService{bs: bs, l: l}
}
