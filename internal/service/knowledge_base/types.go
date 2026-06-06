package knowledge_base

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type knowledgeBaseStorage interface {
	SelectWithoutEmbedding(ctx context.Context) ([]domain.Object, error)
	UpdateEmbeddings(ctx context.Context, id uuid.UUID, embedding []float64) error
	SimilarByEmbedding(ctx context.Context, embedding []float64, limit int) ([]domain.KnowledgeChunk, error)
}

type embedder interface {
	EmbedText(ctx context.Context, text string) ([]float64, error)
}

type llm interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
type aiService struct {
	kbs knowledgeBaseStorage
	emb embedder
	llm llm
	l   *slog.Logger
}

func New(kbs knowledgeBaseStorage, emb embedder, llm llm, l *slog.Logger) *aiService {
	return &aiService{
		kbs: kbs,
		emb: emb,
		l:   l,
		llm: llm,
	}
}
