package reading_sessions

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type readingSessionsStorage interface {
	Save(ctx context.Context, readingSession domain.ReadingSession) error
	AllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.ReadingSession, error)
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type readingSessionsService struct {
	rss readingSessionsStorage
	us  userStorage
	l   *slog.Logger
}

func New(rss readingSessionsStorage, us userStorage, l *slog.Logger) *readingSessionsService {
	return &readingSessionsService{rss: rss, us: us, l: l}
}
