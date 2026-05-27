package recommendation

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type recommendationService interface {
	HomeRecommendation(ctx context.Context, firebaseId string, limit int) ([]domain.BookPreview, error)
	CartRecommendation(ctx context.Context, firebaseId string, limit int) ([]domain.BookPreview, error)
	RecommesPage(ctx context.Context, firebaseId string, limit int) ([]domain.BookPreview, []domain.BookPreview, []domain.BookPreview, error)
	RecommendationByBook(ctx context.Context, firebaseId string, bookId uuid.UUID, limit int) ([]domain.BookPreview, error)
}

type recommendationHandler struct {
	rserv recommendationService
	l     *slog.Logger
}

func New(rserv recommendationService, l *slog.Logger) *recommendationHandler {
	return &recommendationHandler{rserv: rserv, l: l}
}
