package bookviews

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type bookViewsStorage interface {
	SaveOrUpdate(ctx context.Context, userId uuid.UUID, bookId uuid.UUID) error
	LastRecords(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookView, error)
}

type bookStorage interface {
	BookViewPreviews(ctx context.Context, bookViews []domain.BookView) ([]domain.BookViewPreview, error)
}

type historyStorage interface {
	SaveBookView(ctx context.Context, userId uuid.UUID, bookId uuid.UUID) error
	LastBookViews(ctx context.Context, userId uuid.UUID, limit int) ([]domain.BookView, error)
	WarmBookViews(ctx context.Context, userId uuid.UUID, bookViews []domain.BookView, limit int) error
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type bookViewsService struct {
	bvs bookViewsStorage
	bs  bookStorage
	hs  historyStorage
	us  userStorage
	l   *slog.Logger
}

func New(bvs bookViewsStorage, bs bookStorage, hs historyStorage, us userStorage, l *slog.Logger) *bookViewsService {
	return &bookViewsService{bvs: bvs, bs: bs, hs: hs, us: us, l: l}
}
