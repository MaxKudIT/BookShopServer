package physical_books

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type physicalBooksService interface {
	All(ctx context.Context) ([]domain.PhysicalBook, error)
	ById(ctx context.Context, id uuid.UUID) (domain.PhysicalBook, error)
}

type physicalBooksHandler struct {
	pbserv physicalBooksService
	l      *slog.Logger
}

func New(pbserv physicalBooksService, l *slog.Logger) *physicalBooksHandler {
	return &physicalBooksHandler{pbserv: pbserv, l: l}
}
