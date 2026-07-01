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

### Advanced input

```json
{
  "profile": "fast",
  "prompt": "Reply briefly",
  "top_p": 0.9,
  "presence_penalty": 0.1,
  "frequency_penalty": 0.1,
  "seed": 42,
  "response_format": {
    "type": "json_object"
  }
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
- `model` overrides all profile/default model settings.
- `profile` can be `fast`, `smart`, or `local`.
- `temperature`, `max_tokens`, `top_p`, `presence_penalty`, and `frequency_penalty` can be configured globally or per request.
- `seed` and `response_format` are request-only optional fields.
