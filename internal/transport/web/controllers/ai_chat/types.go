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

type aiDialogService interface {
	Ask(ctx context.Context, firebaseId string, chatId uuid.UUID, question string) (domain.AIAskResult, error)
}

type aiChatHandler struct {
	acs aiChatService
	ads aiDialogService
	l   *slog.Logger
}

func New(acs aiChatService, ads aiDialogService, l *slog.Logger) *aiChatHandler {
	return &aiChatHandler{acs: acs, ads: ads, l: l}
}
