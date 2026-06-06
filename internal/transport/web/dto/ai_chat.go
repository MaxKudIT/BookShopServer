package dto

import (
	"time"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

type AIChatDTO struct {
	Title string `json:"title"`
}

type AIMessageDTO struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func AIChatToDomain(id uuid.UUID, aiChatDTO AIChatDTO, createdAt time.Time) domain.AIChat {
	return domain.AIChat{
		Id:        id,
		UserId:    uuid.Nil,
		Title:     aiChatDTO.Title,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

func AIMessageToDomain(id uuid.UUID, chatId uuid.UUID, aiMessageDTO AIMessageDTO, createdAt time.Time) domain.AIMessage {
	return domain.AIMessage{
		Id:        id,
		UserId:    uuid.Nil,
		ChatId:    chatId,
		Role:      aiMessageDTO.Role,
		Content:   aiMessageDTO.Content,
		CreatedAt: createdAt,
	}
}
