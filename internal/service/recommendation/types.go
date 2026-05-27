package recommendation

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type recommendationStorage interface {
	HomeRecommendation(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookPreview, error)
	CartRecommendation(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookPreview, error)
	RecommsPage(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookPreview, []domain.BookPreview, []domain.BookPreview, error)
	RecommendationByBook(ctx context.Context, userId uuid.UUID, bookId uuid.UUID, limit int) ([]domain.BookPreview, error)
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type recommendationService struct {
	rs recommendationStorage
	us userStorage
	l  *slog.Logger
}

func New(rs recommendationStorage, us userStorage, l *slog.Logger) *recommendationService {
	return &recommendationService{rs: rs, us: us, l: l}
}
