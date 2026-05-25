package reading

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type readingService interface {
	Start(ctx context.Context, firebaseId string, bookId uuid.UUID) (domain.ReadingState, error)
	UpdateProgress(ctx context.Context, firebaseId string, bookId uuid.UUID, currentPage int) (domain.ReadingState, error)
	Finish(ctx context.Context, firebaseId string, sessionId uuid.UUID, currentPage int) (domain.ReadingState, error)
}

type readingHandler struct {
	rserv readingService
	l     *slog.Logger
}

func New(rserv readingService, l *slog.Logger) *readingHandler {
	return &readingHandler{rserv: rserv, l: l}
}
