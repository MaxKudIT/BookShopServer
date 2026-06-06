package knowledge_base

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bookshop/internal/domain"
	"github.com/bookshop/internal/utils"
	"github.com/google/uuid"
)

func (nets *netStorage) SelectWithoutEmbedding(ctx context.Context) ([]domain.Object, error) {

	var wearray []domain.Object = make([]domain.Object, 0)
	const SelectWithoutEmbeddingQuery = "SELECT id, content from knowledge_base WHERE embedding IS NULL"
	rows, err := nets.db.QueryContext(ctx, SelectWithoutEmbeddingQuery)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			nets.l.Warn("Query cancelled", "error", err)
			return nil, err
		case errors.Is(err, context.DeadlineExceeded):
			nets.l.Warn("Query timed out", "error", err)
			return nil, err
		default:
			nets.l.Error("Query failed", "error", err)
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var currentObj domain.Object
		if err := rows.Scan(&currentObj.Id, &currentObj.Content); err != nil {
			nets.l.Error("Scan failed", "error", err)
			return nil, err
		}
		wearray = append(wearray, currentObj)
	}

	nets.l.Info("Successfully got records")
	return wearray, nil
}

func (nets *netStorage) UpdateEmbeddings(ctx context.Context, id uuid.UUID, embedding []float64) error {
	const UpdateEmbeddingsQuery = "UPDATE knowledge_base SET embedding = $1::vector, updated_at = NOW() WHERE id = $2"

	if _, err := nets.db.ExecContext(ctx, UpdateEmbeddingsQuery, utils.VectorToSQL(embedding), id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			nets.l.Error("knowledge not found", "error", err)
			return err
		case errors.Is(err, context.Canceled):
			nets.l.Warn("Query cancelled", "error", err)
			return err
		case errors.Is(err, context.DeadlineExceeded):
			nets.l.Warn("Query timed out", "error", err)
			return err
		default:
			nets.l.Error("Query failed", "error", err)
			return err
		}
	}
	nets.l.Info("Successfully updated embedding", "id", id)
	return nil
}

func (nets *netStorage) SimilarByEmbedding(ctx context.Context, embedding []float64, limit int) ([]domain.KnowledgeChunk, error) {
	const query = `
SELECT id, title, content
FROM knowledge_base
WHERE embedding IS NOT NULL
ORDER BY embedding <=> $1::vector
LIMIT $2
`

	rows, err := nets.db.QueryContext(ctx, query, utils.VectorToSQL(embedding), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chunks := make([]domain.KnowledgeChunk, 0, limit)

	for rows.Next() {
		var chunk domain.KnowledgeChunk

		if err := rows.Scan(&chunk.Id, &chunk.Title, &chunk.Content); err != nil {
			return nil, err
		}

		chunks = append(chunks, chunk)
	}

	return chunks, rows.Err()
}
