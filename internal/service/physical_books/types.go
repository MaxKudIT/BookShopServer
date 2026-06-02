package physical_books

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type physicalBooksStorage interface {
	All(ctx context.Context) ([]domain.PhysicalBook, error)
	ById(ctx context.Context, id uuid.UUID) (domain.PhysicalBook, error)
}

type physicalBooksService struct {
	pbs physicalBooksStorage
	l   *slog.Logger
}

func New(pbs physicalBooksStorage, l *slog.Logger) *physicalBooksService {
	return &physicalBooksService{pbs: pbs, l: l}
}
