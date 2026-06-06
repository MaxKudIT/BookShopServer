package ai_chat

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type aiChatService interface {
	CreateChat(ctx context.Context, aiChat domain.AIChat, firebaseId string) (uuid.UUID, error)
	CreateMessage(ctx context.Context, aiMessage domain.AIMessage, firebaseId string) (uuid.UUID, error)
	Messages(ctx context.Context, firebaseId string, chatId uuid.UUID) ([]domain.AIMessage, error)
	DeleteMessages(ctx context.Context, firebaseId string, chatId uuid.UUID) error
}

type aiChatHandler struct {
	acs aiChatService
	l   *slog.Logger
}

func New(acs aiChatService, l *slog.Logger) *aiChatHandler {
	return &aiChatHandler{acs: acs, l: l}
}
