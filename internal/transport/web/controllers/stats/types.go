package stats

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
)

type statsService interface {
	UserStats(ctx context.Context, firebaseId string) (domain.UserStats, error)
}

type statsHandler struct {
	sserv statsService
	l     *slog.Logger
}

func New(sserv statsService, l *slog.Logger) *statsHandler {
	return &statsHandler{sserv: sserv, l: l}
}
