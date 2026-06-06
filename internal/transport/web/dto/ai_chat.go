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

type AIAskDTO struct {
	Question string `json:"question" binding:"required"`
}

type AIAskResponseDTO struct {
	UserMessageId      uuid.UUID `json:"userMessageId"`
	AssistantMessageId uuid.UUID `json:"assistantMessageId"`
	Answer             string    `json:"answer"`
}

type AIChatResponseDTO struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

func AIChatToResponseDTO(aiChat domain.AIChat) AIChatResponseDTO {
	return AIChatResponseDTO{
		Id:        aiChat.Id,
		Title:     aiChat.Title,
		CreatedAt: aiChat.CreatedAt,
		UpdatedAt: aiChat.UpdatedAt,
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
