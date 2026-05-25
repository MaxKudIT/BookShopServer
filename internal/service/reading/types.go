package reading

import (
	"context"
	"log/slog"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type bookStorage interface {
	PagesCountById(ctx context.Context, bookId uuid.UUID) (int, error)
}

type readingStorage interface {
	Start(ctx context.Context, userId uuid.UUID, bookId uuid.UUID, sessionId uuid.UUID, startedAt time.Time) (domain.ReadingState, error)
	UpdateProgress(ctx context.Context, userId uuid.UUID, bookId uuid.UUID, currentPage int, progressPercent int, status domain.Status) (domain.ReadingState, error)
	ActiveSessionBookId(ctx context.Context, userId uuid.UUID, sessionId uuid.UUID) (uuid.UUID, error)
	Finish(ctx context.Context, userId uuid.UUID, sessionId uuid.UUID, currentPage int, progressPercent int, status domain.Status, endedAt time.Time) (domain.ReadingState, error)
}

type readingService struct {
	rs readingStorage
	us userStorage
	bs bookStorage
	l  *slog.Logger
}

func New(rs readingStorage, us userStorage, bs bookStorage, l *slog.Logger) *readingService {
	return &readingService{rs: rs, us: us, bs: bs, l: l}
}
