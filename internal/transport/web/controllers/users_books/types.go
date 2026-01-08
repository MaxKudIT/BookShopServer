package users_books

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
)

type ubService interface {
	Buy(ctx context.Context, firebaseId string, bookId uuid.UUID) error
}

type ubHandler struct {
	ubserv ubService
	l      *slog.Logger
}

func New(ubserv ubService, l *slog.Logger) *ubHandler {
	return &ubHandler{ubserv: ubserv, l: l}
}
