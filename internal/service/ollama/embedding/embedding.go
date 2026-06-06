package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const modelName = "nomic-embed-text"

func (s *Service) EmbedText(ctx context.Context, text string) ([]float64, error) {
	body := embedRequest{
		Model: modelName,
		Input: text,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.client.URL("/api/embed"),
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("ollama embed failed: status %d", res.StatusCode)
	}

	var out embedResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	if len(out.Embeddings) == 0 {
		return nil, errors.New("ollama returned empty embeddings")
	}

	return out.Embeddings[0], nil
}
