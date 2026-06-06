package chat_ai

import (
	"context"
	"log/slog"
)

type aiService interface {
	FillWithoutEmbedding(ctx context.Context) error
	Ask(ctx context.Context, question string) (string, error)
}

type Handler struct {
	service aiService
	l       *slog.Logger
}

func New(service aiService, l *slog.Logger) *Handler {
	return &Handler{
		service: service,
		l:       l,
	}
}
