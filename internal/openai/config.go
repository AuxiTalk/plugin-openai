package openai

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	APIKey           string
	BaseURL          string
	Model            string
	ModelFast        string
	ModelSmart       string
	ModelLocal       string
	Timeout          time.Duration
	Temperature      *float64
	MaxTokens        *int
	TopP             *float64
	PresencePenalty  *float64
	FrequencyPenalty *float64
	Retries          int
	RetryBackoff     time.Duration
	HealthCheck      bool
}

func LoadConfig() (Config, error) {
	cfg := Config{
		APIKey:       os.Getenv("OPENAI_API_KEY"),
		BaseURL:      envOrDefault("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		Model:        envOrDefault("OPENAI_MODEL", "gpt-4o-mini"),
		ModelFast:    os.Getenv("OPENAI_MODEL_FAST"),
		ModelSmart:   os.Getenv("OPENAI_MODEL_SMART"),
		ModelLocal:   os.Getenv("OPENAI_MODEL_LOCAL"),
		Timeout:      60 * time.Second,
		RetryBackoff: time.Second,
	}

	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("OPENAI_API_KEY is required")
	}

	if err := parseDurationEnv("OPENAI_TIMEOUT", &cfg.Timeout); err != nil {
		return Config{}, err
	}
	if err := parseFloatEnv("OPENAI_TEMPERATURE", &cfg.Temperature); err != nil {
		return Config{}, err
	}
	if err := parseIntEnv("OPENAI_MAX_TOKENS", &cfg.MaxTokens); err != nil {
		return Config{}, err
	}
	if err := parseFloatEnv("OPENAI_TOP_P", &cfg.TopP); err != nil {
		return Config{}, err
	}
	if err := parseFloatEnv("OPENAI_PRESENCE_PENALTY", &cfg.PresencePenalty); err != nil {
		return Config{}, err
	}
	if err := parseFloatEnv("OPENAI_FREQUENCY_PENALTY", &cfg.FrequencyPenalty); err != nil {
		return Config{}, err
	}
	if retries := os.Getenv("OPENAI_RETRIES"); retries != "" {
		parsed, err := strconv.Atoi(retries)
		if err != nil {
			return Config{}, fmt.Errorf("invalid OPENAI_RETRIES: %w", err)
		}
		cfg.Retries = parsed
	}
	if err := parseDurationEnv("OPENAI_RETRY_BACKOFF", &cfg.RetryBackoff); err != nil {
		return Config{}, err
	}
	if healthCheck := os.Getenv("OPENAI_HEALTH_CHECK"); healthCheck != "" {
		parsed, err := strconv.ParseBool(healthCheck)
		if err != nil {
			return Config{}, fmt.Errorf("invalid OPENAI_HEALTH_CHECK: %w", err)
		}
		cfg.HealthCheck = parsed
	}

	return cfg, nil
}

func (c Config) ModelForProfile(profile string) string {
	switch profile {
	case "fast":
		if c.ModelFast != "" {
			return c.ModelFast
		}
	case "smart":
		if c.ModelSmart != "" {
			return c.ModelSmart
		}
	case "local":
		if c.ModelLocal != "" {
			return c.ModelLocal
		}
	}
	return c.Model
}

func parseDurationEnv(key string, target *time.Duration) error {
	value := os.Getenv(key)
	if value == "" {
		return nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("invalid %s: %w", key, err)
	}
	*target = parsed
	return nil
}

func parseFloatEnv(key string, target **float64) error {
	value := os.Getenv(key)
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid %s: %w", key, err)
	}
	*target = &parsed
	return nil
}

func parseIntEnv(key string, target **int) error {
	value := os.Getenv(key)
	if value == "" {
		return nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid %s: %w", key, err)
	}
	*target = &parsed
	return nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
