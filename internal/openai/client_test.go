package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestClientCompleteWithSystemAndMessages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if len(req.Messages) != 3 {
			t.Fatalf("expected 3 messages, got %d", len(req.Messages))
		}
		if req.Messages[0].Role != "system" {
			t.Fatalf("expected system first, got %s", req.Messages[0].Role)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"model":"test-model","choices":[{"message":{"role":"assistant","content":"ok"}}]}`))
	}))
	defer server.Close()

	client := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "test-model", Timeout: time.Second})
	_, err := client.Complete(context.Background(), CompleteInput{
		System: "You are helpful.",
		Messages: []InputMessage{
			{Role: "user", Content: "hello"},
		},
		Prompt: "continue",
	})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
}

func TestClientCompleteReturnsAPIErrorBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"bad model"}`, http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "test-model", Timeout: time.Second})
	_, err := client.Complete(context.Background(), CompleteInput{Prompt: "hello"})
	if err == nil || !strings.Contains(err.Error(), "bad model") {
		t.Fatalf("expected api error body, got %v", err)
	}
}

func TestClientCompleteRequiresPromptOrMessages(t *testing.T) {
	client := NewClient(Config{APIKey: "test-key", BaseURL: "http://example.com", Model: "test"})
	if _, err := client.Complete(context.Background(), CompleteInput{}); err == nil {
		t.Fatal("expected prompt or messages required error")
	}
}
