package embedding

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bookshop/internal/service/ollama"
)

func TestEmbedText(t *testing.T) {
	t.Run("sends expected request and returns first embedding", func(t *testing.T) {
		var gotRequest embedRequest
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/embed" {
				t.Fatalf("got path %s", r.URL.Path)
			}
			if r.Method != http.MethodPost {
				t.Fatalf("got method %s", r.Method)
			}
			if err := json.NewDecoder(r.Body).Decode(&gotRequest); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			_ = json.NewEncoder(w).Encode(embedResponse{Embeddings: [][]float64{{1, 2, 3}, {4, 5, 6}}})
		}))
		defer server.Close()

		client := ollama.New(server.URL)
		client.HTTPClient = server.Client()
		got, err := New(client).EmbedText(context.Background(), "text")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotRequest.Model != modelName || gotRequest.Input != "text" {
			t.Fatalf("unexpected request: %+v", gotRequest)
		}
		if len(got) != 3 || got[0] != 1 || got[2] != 3 {
			t.Fatalf("unexpected embedding: %+v", got)
		}
	})

	t.Run("returns error on empty embeddings", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(embedResponse{})
		}))
		defer server.Close()

		client := ollama.New(server.URL)
		client.HTTPClient = server.Client()
		if _, err := New(client).EmbedText(context.Background(), "text"); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("returns error on non success status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "nope", http.StatusBadGateway)
		}))
		defer server.Close()

		client := ollama.New(server.URL)
		client.HTTPClient = server.Client()
		if _, err := New(client).EmbedText(context.Background(), "text"); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("returns error on invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("{"))
		}))
		defer server.Close()

		client := ollama.New(server.URL)
		client.HTTPClient = server.Client()
		if _, err := New(client).EmbedText(context.Background(), "text"); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("returns request error when server is unreachable", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		url := server.URL
		server.Close()

		client := ollama.New(url)
		client.HTTPClient = server.Client()
		if _, err := New(client).EmbedText(context.Background(), "text"); err == nil {
			t.Fatalf("expected error")
		}
	})
}
