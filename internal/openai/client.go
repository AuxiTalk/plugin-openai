package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func (c *Client) Complete(ctx context.Context, input CompleteInput) (CompleteOutput, error) {
	if strings.TrimSpace(input.Prompt) == "" {
		return CompleteOutput{}, fmt.Errorf("prompt is required")
	}

	model := input.Model
	if model == "" {
		model = c.cfg.Model
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
		Model: model,
		Messages: []chatMessage{
			{Role: "user", Content: input.Prompt},
		},
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return CompleteOutput{}, err
	}

	url := strings.TrimRight(c.cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return CompleteOutput{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CompleteOutput{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return CompleteOutput{}, fmt.Errorf("openai-compatible api returned status %d", resp.StatusCode)
	}

	var parsed chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return CompleteOutput{}, err
	}

	if len(parsed.Choices) == 0 {
		return CompleteOutput{}, fmt.Errorf("openai-compatible api returned no choices")
	}

	return CompleteOutput{
		Text:  parsed.Choices[0].Message.Content,
		Model: parsed.Model,
		Usage: parsed.Usage,
	}, nil
}
