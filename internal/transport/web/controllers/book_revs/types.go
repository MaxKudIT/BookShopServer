package book_revs

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type bookRevsService interface {
	Create(ctx context.Context, bookReview domain.BookReview, firebaseId string) (uuid.UUID, error)
	All(ctx context.Context, firebaseId string) ([]domain.BookReview, error)
}

type bookRevsHandler struct {
	brserv bookRevsService
	l      *slog.Logger
}

func New(brserv bookRevsService, l *slog.Logger) *bookRevsHandler {
	return &bookRevsHandler{brserv: brserv, l: l}
}
