package domain

import (
	"time"

	"github.com/google/uuid"
)

type KnowledgeBase struct {
	Id         uuid.UUID
	Content    string
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	Embedding  *[][]float64
	Title      string
	SourceType string
	BookId     *uuid.UUID
}

type Object struct {
	Id      uuid.UUID
	Content string
}

type KnowledgeChunk struct {
	Id      uuid.UUID
	Content string
	Title   string
}
