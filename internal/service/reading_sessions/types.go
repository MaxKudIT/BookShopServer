package reading_sessions

import (
	"context"
	"log/slog"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type readingSessionsStorage interface {
	Save(ctx context.Context, readingSession domain.ReadingSession) error
	Close(ctx context.Context, userId uuid.UUID, sessionId uuid.UUID, endedAt time.Time) (domain.ReadingSession, error)
	AllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.ReadingSession, error)
	LastReadingBookRecords(ctx context.Context, userId uuid.UUID, limit int) ([]domain.LastReadingBook, error)
}

type bookStorage interface {
	ReadingBookPreviews(ctx context.Context, userId uuid.UUID, lastReadingBooks []domain.LastReadingBook) ([]domain.ReadingBookPreview, error)
}

type historyStorage interface {
	LastReadingBooks(ctx context.Context, userId uuid.UUID, limit int) ([]domain.LastReadingBook, error)
	WarmReadingBooks(ctx context.Context, userId uuid.UUID, lastReadingBooks []domain.LastReadingBook, limit int) error
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type readingSessionsService struct {
	rss readingSessionsStorage
	bs  bookStorage
	hs  historyStorage
	us  userStorage
	l   *slog.Logger
}

func New(rss readingSessionsStorage, bs bookStorage, hs historyStorage, us userStorage, l *slog.Logger) *readingSessionsService {
	return &readingSessionsService{rss: rss, bs: bs, hs: hs, us: us, l: l}
}
