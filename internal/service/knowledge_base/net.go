package knowledge_base

import (
	"context"
	"fmt"
	"github.com/bookshop/internal/domain"
	"strings"
)

func (s *aiService) FillWithoutEmbedding(ctx context.Context) error {

	items, err := s.kbs.SelectWithoutEmbedding(ctx)
	if err != nil {
		s.l.Error("failed to get knowledge without embedding", "error", err)
		return err
	}

	for _, item := range items {
		embedding, err := s.emb.EmbedText(ctx, item.Content)
		if err != nil {
			s.l.Error("failed to embed knowledge", "id", item.Id, "error", err)
			return err
		}

		if err := s.kbs.UpdateEmbeddings(ctx, item.Id, embedding); err != nil {
			s.l.Error("failed to update knowledge embedding", "id", item.Id, "error", err)
			return err
		}
	}

	s.l.Info("successfully filled knowledge embeddings", "count", len(items))
	return nil
}

func (s *aiService) Ask(ctx context.Context, question string) (string, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return "", fmt.Errorf("question is empty")
	}

	embedding, err := s.emb.EmbedText(ctx, question)
	if err != nil {
		return "", err
	}

	chunks, err := s.kbs.SimilarByEmbedding(ctx, embedding, 5)
	if err != nil {
		return "", err
	}

	prompt := buildPrompt(question, chunks)

	answer, err := s.llm.Generate(ctx, prompt)
	if err != nil {
		return "", err
	}

	return answer, nil
}

func buildPrompt(question string, chunks []domain.KnowledgeChunk) string {
	var b strings.Builder

	b.WriteString("Ты помощник онлайн-библиотеки BookShop.\n")
	b.WriteString("Отвечай только на основе контекста ниже.\n")
	b.WriteString("Если в контексте нет ответа, скажи, что данных недостаточно.\n\n")

	b.WriteString("Контекст:\n")
	for i, chunk := range chunks {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, chunk.Title))
		b.WriteString(chunk.Content)
		b.WriteString("\n\n")
	}

	b.WriteString("Вопрос пользователя:\n")
	b.WriteString(question)

	return b.String()
}
