package openai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientComplete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("missing authorization header")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"model":"test-model",
			"choices":[{"message":{"role":"assistant","content":"hello from test"}}],
			"usage":{"prompt_tokens":5,"completion_tokens":3,"total_tokens":8}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Model:   "test-model",
		Timeout: time.Second,
	})

	out, err := client.Complete(context.Background(), CompleteInput{Prompt: "hello"})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if out.Text != "hello from test" {
		t.Fatalf("unexpected text: %s", out.Text)
	}
	if out.Usage.TotalTokens != 8 {
		t.Fatalf("unexpected usage: %+v", out.Usage)
	}
}

func TestClientCompleteRequiresPrompt(t *testing.T) {
	client := NewClient(Config{APIKey: "test-key", BaseURL: "http://example.com", Model: "test"})
	if _, err := client.Complete(context.Background(), CompleteInput{}); err == nil {
		t.Fatal("expected prompt required error")
	}
}
