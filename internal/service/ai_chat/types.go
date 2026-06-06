package ai_chat

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type aiChatStorage interface {
	SaveChat(ctx context.Context, aiChat domain.AIChat) error
	SaveMessage(ctx context.Context, aiMessage domain.AIMessage) error
	ChatByUserId(ctx context.Context, userId uuid.UUID) (domain.AIChat, error)
	MessagesByChatId(ctx context.Context, userId uuid.UUID, chatId uuid.UUID) ([]domain.AIMessage, error)
	DeleteMessagesByChatId(ctx context.Context, userId uuid.UUID, chatId uuid.UUID) error
}

type userStorage interface {
	UserByFirebaseId(ctx context.Context, firebaseId string) (uuid.UUID, error)
}

type aiChatService struct {
	acs aiChatStorage
	us  userStorage
	l   *slog.Logger
}

func New(acs aiChatStorage, us userStorage, l *slog.Logger) *aiChatService {
	return &aiChatService{acs: acs, us: us, l: l}
}
