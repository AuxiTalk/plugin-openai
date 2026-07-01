package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	cfg        Config
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *Client) Health(ctx context.Context) error {
	if !c.cfg.HealthCheck {
		return nil
	}

	url := strings.TrimRight(c.cfg.BaseURL, "/") + "/models"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("health check failed status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return nil
}

func (c *Client) Complete(ctx context.Context, input CompleteInput) (CompleteOutput, error) {
	messages, err := buildMessages(input)
	if err != nil {
		return CompleteOutput{}, err
	}

	model := input.Model
	if model == "" {
		model = c.cfg.ModelForProfile(input.Profile)
	}

	temperature := input.Temperature
	if temperature == nil {
		temperature = c.cfg.Temperature
	}

	maxTokens := input.MaxTokens
	if maxTokens == nil {
		maxTokens = c.cfg.MaxTokens
	}

	body := chatRequest{
		Model:            model,
		Messages:         messages,
		Temperature:      temperature,
		MaxTokens:        maxTokens,
		TopP:             firstFloat(input.TopP, c.cfg.TopP),
		PresencePenalty:  firstFloat(input.PresencePenalty, c.cfg.PresencePenalty),
		FrequencyPenalty: firstFloat(input.FrequencyPenalty, c.cfg.FrequencyPenalty),
		Seed:             input.Seed,
		ResponseFormat:   input.ResponseFormat,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return CompleteOutput{}, err
	}

	return c.doComplete(ctx, payload)
}

func (c *Client) doComplete(ctx context.Context, payload []byte) (CompleteOutput, error) {
	var lastErr error
	attempts := c.cfg.Retries + 1
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return CompleteOutput{}, ctx.Err()
			case <-time.After(c.cfg.RetryBackoff * time.Duration(attempt)):
			}
		}

		result, retry, err := c.doCompleteOnce(ctx, payload)
		if err == nil {
			return result, nil
		}
		lastErr = err
		if !retry {
			break
		}
	}
	return CompleteOutput{}, lastErr
}

func (c *Client) doCompleteOnce(ctx context.Context, payload []byte) (CompleteOutput, bool, error) {
	url := strings.TrimRight(c.cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return CompleteOutput{}, false, err
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CompleteOutput{}, true, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		err := fmt.Errorf("openai-compatible api returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
		return CompleteOutput{}, isRetryableStatus(resp.StatusCode), err
	}

	var parsed chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return CompleteOutput{}, false, err
	}

	if len(parsed.Choices) == 0 {
		return CompleteOutput{}, false, fmt.Errorf("openai-compatible api returned no choices")
	}

	return CompleteOutput{
		Text:  parsed.Choices[0].Message.Content,
		Model: parsed.Model,
		Usage: parsed.Usage,
	}, false, nil
}

func buildMessages(input CompleteInput) ([]chatMessage, error) {
	messages := make([]chatMessage, 0, len(input.Messages)+2)

	if strings.TrimSpace(input.System) != "" {
		messages = append(messages, chatMessage{Role: "system", Content: input.System})
	}

	for _, msg := range input.Messages {
		if strings.TrimSpace(msg.Role) == "" || strings.TrimSpace(msg.Content) == "" {
			return nil, fmt.Errorf("messages require role and content")
		}
		messages = append(messages, chatMessage{Role: msg.Role, Content: msg.Content})
	}

	if strings.TrimSpace(input.Prompt) != "" {
		messages = append(messages, chatMessage{Role: "user", Content: input.Prompt})
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("prompt or messages are required")
	}

	return messages, nil
}

func firstFloat(value, fallback *float64) *float64 {
	if value != nil {
		return value
	}
	return fallback
}

func isRetryableStatus(status int) bool {
	switch status {
	case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}
