# Capabilities

## `ai.complete`

Generates a completion using an OpenAI-compatible chat completion API.

### Simple input

```json
{
  "prompt": "Write a short reply",
  "model": "optional-model",
  "temperature": 0.7,
  "max_tokens": 1024
}
```

### Chat input

```json
{
  "system": "You are a helpful assistant.",
  "messages": [
    {
      "role": "user",
      "content": "Hello"
    }
  ],
  "prompt": "Continue the conversation"
}
```

### Output

```json
{
  "text": "Generated response",
  "model": "model-name",
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 20,
    "total_tokens": 30
  }
}
```

### Notes

- Either `prompt` or `messages` is required.
- `system` is added as the first message when provided.
- `prompt` is appended as a final user message when provided with `messages`.
- `model` defaults to `OPENAI_MODEL`.
- `temperature` defaults to `OPENAI_TEMPERATURE` when set.
- `max_tokens` defaults to `OPENAI_MAX_TOKENS` when set.
