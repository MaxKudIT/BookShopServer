package fav

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type favService interface {
	Create(ctx context.Context, fav domain.Fav, firebaseId string) (uuid.UUID, error)
}

type favHandler struct {
	fserv favService
	l     *slog.Logger
}

func New(fserv favService, l *slog.Logger) *favHandler {
	return &favHandler{fserv: fserv, l: l}
}
