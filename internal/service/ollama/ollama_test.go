package ollama

import "testing"

func TestClientURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		path    string
		want    string
	}{
		{"trims base slash and path slash", "http://localhost:11434/", "/api/generate", "http://localhost:11434/api/generate"},
		{"keeps base without slash", "http://localhost:11434", "/api/embed", "http://localhost:11434/api/embed"},
		{"accepts path without leading slash", "http://localhost:11434", "api/tags", "http://localhost:11434/api/tags"},
		{"keeps nested base path", "http://localhost:11434/ollama/", "/api/generate", "http://localhost:11434/ollama/api/generate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(tt.baseURL)
			if got := client.URL(tt.path); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
