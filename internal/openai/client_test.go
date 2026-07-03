package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
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

	client := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "test-model", Timeout: time.Second})

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

func TestClientCompleteUsesProfileAndAdvancedParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Model != "fast-model" {
			t.Fatalf("expected fast-model, got %s", req.Model)
		}
		if req.TopP == nil || *req.TopP != 0.9 {
			t.Fatalf("expected top_p 0.9")
		}
		if req.Seed == nil || *req.Seed != 42 {
			t.Fatalf("expected seed 42")
		}
		w.Write([]byte(`{"model":"fast-model","choices":[{"message":{"role":"assistant","content":"ok"}}]}`))
	}))
	defer server.Close()

	topP := 0.9
	seed := 42
	client := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "default-model", ModelFast: "fast-model", Timeout: time.Second})
	_, err := client.Complete(context.Background(), CompleteInput{Prompt: "hello", Profile: "fast", TopP: &topP, Seed: &seed})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
}

func TestClientCompleteRetriesRetryableStatus(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&calls, 1)
		if count == 1 {
			http.Error(w, "temporary", http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte(`{"model":"test-model","choices":[{"message":{"role":"assistant","content":"ok"}}]}`))
	}))
	defer server.Close()

	client := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "test-model", Timeout: time.Second, Retries: 1, RetryBackoff: time.Millisecond})
	_, err := client.Complete(context.Background(), CompleteInput{Prompt: "hello"})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestClientHealthOptional(t *testing.T) {
	client := NewClient(Config{})
	if err := client.Health(context.Background()); err != nil {
		t.Fatalf("health disabled should pass: %v", err)
	}
}

func TestClientHealthEnabled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()

	client := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Timeout: time.Second, HealthCheck: true})
	if err := client.Health(context.Background()); err != nil {
		t.Fatalf("health: %v", err)
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

func TestClientCompleteReturnsNoChoicesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"model":"test-model","choices":[]}`))
	}))
	defer server.Close()

	client := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "test-model", Timeout: time.Second})
	_, err := client.Complete(context.Background(), CompleteInput{Prompt: "hello"})
	if err == nil || !strings.Contains(err.Error(), "no choices") {
		t.Fatalf("expected no choices error, got %v", err)
	}
}

func TestClientCompleteDoesNotRetryBadRequest(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "test-model", Timeout: time.Second, Retries: 3, RetryBackoff: time.Millisecond})
	_, err := client.Complete(context.Background(), CompleteInput{Prompt: "hello"})
	if err == nil {
		t.Fatal("expected bad request error")
	}
	if calls != 1 {
		t.Fatalf("expected no retry on bad request, got %d calls", calls)
	}
}

func TestClientHealthReturnsErrorBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "models unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Timeout: time.Second, HealthCheck: true})
	err := client.Health(context.Background())
	if err == nil || !strings.Contains(err.Error(), "models unavailable") {
		t.Fatalf("expected health error body, got %v", err)
	}
}
