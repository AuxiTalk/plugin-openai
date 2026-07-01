package openai

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	APIKey      string
	BaseURL     string
	Model       string
	Timeout     time.Duration
	Temperature *float64
	MaxTokens   *int
}

func LoadConfig() (Config, error) {
	cfg := Config{
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		BaseURL: envOrDefault("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		Model:   envOrDefault("OPENAI_MODEL", "gpt-4o-mini"),
		Timeout: 60 * time.Second,
	}

	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("OPENAI_API_KEY is required")
	}

	if timeout := os.Getenv("OPENAI_TIMEOUT"); timeout != "" {
		parsed, err := time.ParseDuration(timeout)
		if err != nil {
			return Config{}, fmt.Errorf("invalid OPENAI_TIMEOUT: %w", err)
		}
		cfg.Timeout = parsed
	}

	if temperature := os.Getenv("OPENAI_TEMPERATURE"); temperature != "" {
		parsed, err := strconv.ParseFloat(temperature, 64)
		if err != nil {
			return Config{}, fmt.Errorf("invalid OPENAI_TEMPERATURE: %w", err)
		}
		cfg.Temperature = &parsed
	}

	if maxTokens := os.Getenv("OPENAI_MAX_TOKENS"); maxTokens != "" {
		parsed, err := strconv.Atoi(maxTokens)
		if err != nil {
			return Config{}, fmt.Errorf("invalid OPENAI_MAX_TOKENS: %w", err)
		}
		cfg.MaxTokens = &parsed
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
