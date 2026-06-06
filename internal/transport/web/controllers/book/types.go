package book

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type bookService interface {
	AllBooks(ctx context.Context, firebaseId string) ([]domain.BookPreview, error)
	AllMyBooks(ctx context.Context, firebaseId string) ([]domain.BookPreview, error)
	AllNotMyBooks(ctx context.Context, firebaseId string) ([]domain.BookPreview, error)
	BookById(ctx context.Context, firebaseId string, bookId uuid.UUID) (domain.Book, error)
	IsMyBook(ctx context.Context, firebaseId string, bookId uuid.UUID) (bool, error)
}

type bookHandler struct {
	bs bookService
	l  *slog.Logger
}

func New(bs bookService, l *slog.Logger) *bookHandler {
	return &bookHandler{bs: bs, l: l}
}
