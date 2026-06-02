package bookviews

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type bookViewsService interface {
	SaveOrUpdate(ctx context.Context, firebaseId string, bookId uuid.UUID) error
	LastRecords(ctx context.Context, firebaseId string, limit int) ([]domain.BookViewPreview, error)
}

type bookViewsHandler struct {
	bvserv bookViewsService
	l      *slog.Logger
}

func New(bvserv bookViewsService, l *slog.Logger) *bookViewsHandler {
	return &bookViewsHandler{bvserv: bvserv, l: l}
}
