package book

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type bookService interface {
	AllBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error)
	AllMyBooks(ctx context.Context, userId uuid.UUID) ([]domain.BookPreview, error)
	BookById(ctx context.Context, id uuid.UUID) (domain.Book, error)
}

type bookHandler struct {
	bs bookService
	l  *slog.Logger
}

func New(bs bookService, l *slog.Logger) *bookHandler {
	return &bookHandler{bs: bs, l: l}
}
