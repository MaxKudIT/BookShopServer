package llm

import (
	"github.com/bookshop/internal/service/ollama"
)

type Service struct {
	client *ollama.Client
}

func New(client *ollama.Client) *Service {
	return &Service{client: client}
}

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type generateResponse struct {
	Response string `json:"response"`
}
