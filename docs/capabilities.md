# Capabilities

## `ai.complete`

Generates a completion using an OpenAI-compatible chat completion API.

### Input

```json
{
  "prompt": "Write a short reply",
  "model": "optional-model",
  "temperature": 0.7,
  "max_tokens": 1024
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

- `prompt` is required.
- `model` defaults to `OPENAI_MODEL`.
- `temperature` defaults to `OPENAI_TEMPERATURE` when set.
- `max_tokens` defaults to `OPENAI_MAX_TOKENS` when set.
