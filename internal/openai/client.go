package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	messages, err := buildMessages(input)
	if err != nil {
		return CompleteOutput{}, err
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
		Model:       model,
		Messages:    messages,
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
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return CompleteOutput{}, fmt.Errorf("openai-compatible api returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
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
