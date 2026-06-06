package ai_chat

import (
	"context"
	"errors"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (acs *aiChatStorage) SaveChat(ctx context.Context, aiChat domain.AIChat) error {
	const CreateAIChatQuery = `
		INSERT INTO ai_chats (id, user_id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	if _, err := acs.db.ExecContext(
		ctx,
		CreateAIChatQuery,
		aiChat.Id,
		aiChat.UserId,
		aiChat.Title,
		aiChat.CreatedAt,
		aiChat.UpdatedAt,
	); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			acs.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			acs.l.Warn("Query timed out", "error", err)
			return err
		default:
			acs.l.Error("Query failed", "error", err)
			return err
		}
	}

	acs.l.Info("Successfully saved ai chat", "id", aiChat.Id)
	return nil
}

func (acs *aiChatStorage) SaveMessage(ctx context.Context, aiMessage domain.AIMessage) error {
	const CreateAIMessageQuery = `
		INSERT INTO ai_messages (id, user_id, chat_id, role, content, created_at)
		SELECT $1, $2, $3, $4, $5, $6
		WHERE EXISTS (
			SELECT 1
			FROM ai_chats
			WHERE id = $3 AND user_id = $2
		)
	`

	result, err := acs.db.ExecContext(
		ctx,
		CreateAIMessageQuery,
		aiMessage.Id,
		aiMessage.UserId,
		aiMessage.ChatId,
		aiMessage.Role,
		aiMessage.Content,
		aiMessage.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			acs.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			acs.l.Warn("Query timed out", "error", err)
			return err
		default:
			acs.l.Error("Query failed", "error", err)
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		acs.l.Error("Rows affected failed", "error", err)
		return err
	}
	if rowsAffected == 0 {
		acs.l.Error("ai chat not found for user", "chatId", aiMessage.ChatId, "userId", aiMessage.UserId)
		return errors.New("ai chat not found")
	}

	acs.l.Info("Successfully saved ai message", "id", aiMessage.Id)
	return nil
}

func (acs *aiChatStorage) MessagesByChatId(ctx context.Context, userId uuid.UUID, chatId uuid.UUID) ([]domain.AIMessage, error) {
	const GetAIMessagesQuery = `
		SELECT m.id, m.user_id, m.chat_id, m.role, m.content, m.created_at
		FROM ai_messages m
		INNER JOIN ai_chats c ON c.id = m.chat_id
		WHERE m.chat_id = $1 AND c.user_id = $2
		ORDER BY m.created_at ASC
	`

	rows, err := acs.db.QueryContext(ctx, GetAIMessagesQuery, chatId, userId)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			acs.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			acs.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			acs.l.Error("Query failed", "error", err)
			return nil, err
		}
	}
	defer rows.Close()

	aiMessages := make([]domain.AIMessage, 0)
	for rows.Next() {
		var aiMessage domain.AIMessage
		if err := rows.Scan(
			&aiMessage.Id,
			&aiMessage.UserId,
			&aiMessage.ChatId,
			&aiMessage.Role,
			&aiMessage.Content,
			&aiMessage.CreatedAt,
		); err != nil {
			acs.l.Error("Scan failed", "error", err)
			return nil, err
		}
		aiMessages = append(aiMessages, aiMessage)
	}

	if err := rows.Err(); err != nil {
		acs.l.Error("Rows failed", "error", err)
		return nil, err
	}

	acs.l.Info("Successfully got ai messages", "chatId", chatId)
	return aiMessages, nil
}

func (acs *aiChatStorage) DeleteMessagesByChatId(ctx context.Context, userId uuid.UUID, chatId uuid.UUID) error {
	const DeleteAIMessagesQuery = `
		DELETE FROM ai_messages
		WHERE chat_id = $1
			AND user_id = $2
	`

	if _, err := acs.db.ExecContext(ctx, DeleteAIMessagesQuery, chatId, userId); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			acs.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			acs.l.Warn("Query timed out", "error", err)
			return err
		default:
			acs.l.Error("Query failed", "error", err)
			return err
		}
	}

	acs.l.Info("Successfully deleted ai messages", "chatId", chatId)
	return nil
}
