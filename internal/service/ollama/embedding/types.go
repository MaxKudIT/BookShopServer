package embedding

import "github.com/bookshop/internal/service/ollama"

type Service struct {
	client *ollama.Client
}

func New(client *ollama.Client) *Service {
	return &Service{
		client: client,
	}
}

type embedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embedResponse struct {
	Embeddings [][]float64 `json:"embeddings"`
}
