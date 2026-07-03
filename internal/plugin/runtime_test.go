package plugin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/auxitalk/plugin-openai/internal/openai"
)

func TestRuntimeCapabilityCallCompletesWithFakeServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"model":"test-model","choices":[{"message":{"role":"assistant","content":"runtime ok"}}]}`))
	}))
	defer server.Close()

	var output bytes.Buffer
	var logs bytes.Buffer
	runtime := NewRuntime(strings.NewReader(""), &output, &logs, openai.Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Model:   "test-model",
		Timeout: time.Second,
	})

	input, _ := json.Marshal(openai.CompleteInput{Prompt: "hello"})
	params, _ := json.Marshal(struct {
		Name  string          `json:"name"`
		Input json.RawMessage `json:"input"`
	}{Name: "ai.complete", Input: input})
	result, err := runtime.capabilityCall(params)
	if err != nil {
		t.Fatalf("capability call: %v", err)
	}
	out, ok := result.(openai.CompleteOutput)
	if !ok {
		t.Fatalf("unexpected result type: %T", result)
	}
	if out.Text != "runtime ok" {
		t.Fatalf("unexpected text: %s", out.Text)
	}
	if strings.Contains(logs.String(), "test-key") {
		t.Fatal("logs must not contain API key")
	}
}

func TestRuntimeRejectsUnknownCapability(t *testing.T) {
	runtime := NewRuntime(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{}, openai.Config{APIKey: "test-key", Timeout: time.Second})
	params, _ := json.Marshal(map[string]any{"name": "ai.unknown", "input": map[string]any{}})
	_, err := runtime.capabilityCall(params)
	if err == nil || !strings.Contains(err.Error(), "capability not found") {
		t.Fatalf("expected capability not found, got %v", err)
	}
}
