package ai_chat

import (
	"context"
	"errors"
	"strings"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (acs *aiChatService) CreateChat(ctx context.Context, aiChat domain.AIChat, firebaseId string) (uuid.UUID, error) {
	userId, err := acs.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		acs.l.Error("Error getting userId by firebaseId", "error", err)
		return uuid.Nil, err
	}

	aiChat.UserId = userId
	aiChat.Title = strings.TrimSpace(aiChat.Title)
	if aiChat.Title == "" {
		aiChat.Title = "Новый чат"
	}

	if err := acs.acs.SaveChat(ctx, aiChat); err != nil {
		acs.l.Error("Error saving ai chat", "error", err)
		return aiChat.Id, err
	}

	acs.l.Info("Successfully created ai chat", "id", aiChat.Id)
	return aiChat.Id, nil
}

func (acs *aiChatService) CreateMessage(ctx context.Context, aiMessage domain.AIMessage, firebaseId string) (uuid.UUID, error) {
	userId, err := acs.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		acs.l.Error("Error getting userId by firebaseId", "error", err)
		return uuid.Nil, err
	}

	aiMessage.UserId = userId
	aiMessage.Role = strings.TrimSpace(aiMessage.Role)
	aiMessage.Content = strings.TrimSpace(aiMessage.Content)

	if aiMessage.Role != "user" && aiMessage.Role != "assistant" {
		return uuid.Nil, errors.New("invalid ai message role")
	}
	if aiMessage.Content == "" {
		return uuid.Nil, errors.New("ai message content is empty")
	}

	if err := acs.acs.SaveMessage(ctx, aiMessage); err != nil {
		acs.l.Error("Error saving ai message", "error", err)
		return aiMessage.Id, err
	}

	acs.l.Info("Successfully created ai message", "id", aiMessage.Id)
	return aiMessage.Id, nil
}

func (acs *aiChatService) CurrentChat(ctx context.Context, firebaseId string) (domain.AIChat, error) {
	userId, err := acs.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		acs.l.Error("Error getting userId by firebaseId", "error", err)
		return domain.AIChat{}, err
	}

	aiChat, err := acs.acs.ChatByUserId(ctx, userId)
	if err != nil {
		acs.l.Error("Error getting current ai chat", "error", err)
		return domain.AIChat{}, err
	}

	acs.l.Info("Successfully got current ai chat", "id", aiChat.Id)
	return aiChat, nil
}

func (acs *aiChatService) Messages(ctx context.Context, firebaseId string, chatId uuid.UUID) ([]domain.AIMessage, error) {
	userId, err := acs.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		acs.l.Error("Error getting userId by firebaseId", "error", err)
		return nil, err
	}

	aiMessages, err := acs.acs.MessagesByChatId(ctx, userId, chatId)
	if err != nil {
		acs.l.Error("Error getting ai messages", "error", err)
		return nil, err
	}

	acs.l.Info("Successfully got ai messages", "chatId", chatId)
	return aiMessages, nil
}

func (acs *aiChatService) DeleteMessages(ctx context.Context, firebaseId string, chatId uuid.UUID) error {
	userId, err := acs.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		acs.l.Error("Error getting userId by firebaseId", "error", err)
		return err
	}

	if err := acs.acs.DeleteMessagesByChatId(ctx, userId, chatId); err != nil {
		acs.l.Error("Error deleting ai messages", "error", err)
		return err
	}

	acs.l.Info("Successfully deleted ai messages", "chatId", chatId)
	return nil
}
