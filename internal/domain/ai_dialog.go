package domain

import "github.com/google/uuid"

type AIAskResult struct {
	UserMessageId      uuid.UUID
	AssistantMessageId uuid.UUID
	Answer             string
}
