package ai_dialog

import (
	"context"
	"log/slog"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type aiChatService interface {
	CreateMessage(ctx context.Context, aiMessage domain.AIMessage, firebaseId string) (uuid.UUID, error)
}

type knowledgeBaseService interface {
	Ask(ctx context.Context, question string) (string, error)
}

type aiDialogService struct {
	acs aiChatService
	kbs knowledgeBaseService
	l   *slog.Logger
}

func New(acs aiChatService, kbs knowledgeBaseService, l *slog.Logger) *aiDialogService {
	return &aiDialogService{
		acs: acs,
		kbs: kbs,
		l:   l,
	}
}
