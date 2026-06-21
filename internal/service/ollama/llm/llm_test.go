package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bookshop/internal/service/ollama"
)

func TestGenerate(t *testing.T) {
	t.Run("sends expected request and returns response", func(t *testing.T) {
		var gotRequest generateRequest
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/generate" {
				t.Fatalf("got path %s", r.URL.Path)
			}
			if r.Method != http.MethodPost {
				t.Fatalf("got method %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Fatalf("missing json content type")
			}
			if err := json.NewDecoder(r.Body).Decode(&gotRequest); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			_ = json.NewEncoder(w).Encode(generateResponse{Response: "hello"})
		}))
		defer server.Close()

		client := ollama.New(server.URL)
		client.HTTPClient = server.Client()
		got, err := New(client).Generate(context.Background(), "prompt")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "hello" {
			t.Fatalf("got %q", got)
		}
		if gotRequest.Model != modelName || gotRequest.Prompt != "prompt" || gotRequest.Stream {
			t.Fatalf("unexpected request: %+v", gotRequest)
		}
	})

	t.Run("returns error on non success status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "nope", http.StatusInternalServerError)
		}))
		defer server.Close()

		client := ollama.New(server.URL)
		client.HTTPClient = server.Client()
		if _, err := New(client).Generate(context.Background(), "prompt"); err == nil {
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
		if _, err := New(client).Generate(context.Background(), "prompt"); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("returns request error when server is unreachable", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		url := server.URL
		server.Close()

		client := ollama.New(url)
		client.HTTPClient = server.Client()
		if _, err := New(client).Generate(context.Background(), "prompt"); err == nil {
			t.Fatalf("expected error")
		}
	})
}
