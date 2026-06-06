package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const modelName = "llama3.2:3b"

func (s *Service) Generate(ctx context.Context, prompt string) (string, error) {
	body := generateRequest{
		Model:  modelName,
		Prompt: prompt,
		Stream: false,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.client.URL("/api/generate"),
		bytes.NewReader(b),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("ollama generate failed: status %d", res.StatusCode)
	}

	var out generateResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", err
	}

	return out.Response, nil
}
