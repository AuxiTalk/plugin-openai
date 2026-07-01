package openai

import (
	"testing"
)

func TestLoadConfigRequiresKey(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected error when API key is missing")
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-key")
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Model != "gpt-4o-mini" {
		t.Fatalf("expected default model gpt-4o-mini, got %s", cfg.Model)
	}
}
