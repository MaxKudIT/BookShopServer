package ai_dialog

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (ads *aiDialogService) Ask(ctx context.Context, firebaseId string, chatId uuid.UUID, question string) (domain.AIAskResult, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return domain.AIAskResult{}, fmt.Errorf("question is empty")
	}

	userMessage := domain.AIMessage{
		Id:        uuid.New(),
		ChatId:    chatId,
		Role:      "user",
		Content:   question,
		CreatedAt: time.Now(),
	}

	userMessageId, err := ads.acs.CreateMessage(ctx, userMessage, firebaseId)
	if err != nil {
		ads.l.Error("failed to save user ai message", "error", err)
		return domain.AIAskResult{}, err
	}

	answer, err := ads.kbs.Ask(ctx, question)
	if err != nil {
		ads.l.Error("failed to ask knowledge base", "error", err)
		return domain.AIAskResult{}, err
	}

	assistantMessage := domain.AIMessage{
		Id:        uuid.New(),
		ChatId:    chatId,
		Role:      "assistant",
		Content:   answer,
		CreatedAt: time.Now(),
	}

	assistantMessageId, err := ads.acs.CreateMessage(ctx, assistantMessage, firebaseId)
	if err != nil {
		ads.l.Error("failed to save assistant ai message", "error", err)
		return domain.AIAskResult{}, err
	}

	ads.l.Info("successfully completed ai dialog ask", "chatId", chatId)
	return domain.AIAskResult{
		UserMessageId:      userMessageId,
		AssistantMessageId: assistantMessageId,
		Answer:             answer,
	}, nil
}
