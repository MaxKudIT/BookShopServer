package reading_sessions

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type readingSessionsService interface {
	Create(ctx context.Context, readingSession domain.ReadingSession, firebaseId string) (uuid.UUID, error)
	Close(ctx context.Context, firebaseId string, sessionId uuid.UUID) (domain.ReadingSession, error)
	All(ctx context.Context, firebaseId string) ([]domain.ReadingSession, error)
	LastReadingBooks(ctx context.Context, firebaseId string, limit int) ([]domain.ReadingBookPreview, error)
}

type readingSessionsHandler struct {
	rsserv readingSessionsService
	l      *slog.Logger
}

func New(rsserv readingSessionsService, l *slog.Logger) *readingSessionsHandler {
	return &readingSessionsHandler{rsserv: rsserv, l: l}
}
