package openai

import "testing"

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
	if cfg.Retries != 0 {
		t.Fatalf("expected retries disabled by default")
	}
	if cfg.HealthCheck {
		t.Fatalf("expected health check disabled by default")
	}
}

func TestLoadConfigOptionalValues(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("OPENAI_RETRIES", "2")
	t.Setenv("OPENAI_HEALTH_CHECK", "true")
	t.Setenv("OPENAI_MODEL_FAST", "fast-model")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Retries != 2 {
		t.Fatalf("expected retries 2, got %d", cfg.Retries)
	}
	if !cfg.HealthCheck {
		t.Fatalf("expected health check enabled")
	}
	if cfg.ModelForProfile("fast") != "fast-model" {
		t.Fatalf("expected fast profile model")
	}
}
