package page

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type pageStorage interface {
	AllPagesOfBook(ctx context.Context, bookId uuid.UUID) (int, error)
	PageByNumber(ctx context.Context, pageNumber int, bookId uuid.UUID) (domain.Page, error)
}

type pageService struct {
	ps pageStorage
	l  *slog.Logger
}

func New(ps pageStorage, l *slog.Logger) *pageService {
	return &pageService{ps: ps, l: l}
}
