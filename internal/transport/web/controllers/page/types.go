package page

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type pageService interface {
	AllPagesOfBook(ctx context.Context, bookId uuid.UUID) (int, error)
	PageByNumber(ctx context.Context, pageNumber int, bookId uuid.UUID) (domain.Page, error)
}

type pageHandler struct {
	ps pageService
	l  *slog.Logger
}

func New(ps pageService, l *slog.Logger) *pageHandler {
	return &pageHandler{ps: ps, l: l}
}
