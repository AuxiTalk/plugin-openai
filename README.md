# AuxiTalk OpenAI Plugin

OpenAI-compatible AI plugin for AuxiTalk.

> Portuguese documentation: [README.pt-BR.md](README.pt-BR.md)

## Overview

This plugin exposes the `ai.complete` capability and connects AuxiTalk to OpenAI-compatible chat completion APIs.

Compatible providers include:

- OpenAI
- Groq
- Together.ai
- Ollama
- LM Studio
- any OpenAI-compatible `/chat/completions` API

## Configuration

Use environment variables:

```txt
OPENAI_API_KEY
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-4o-mini
OPENAI_TIMEOUT=60s
OPENAI_TEMPERATURE=0.7
OPENAI_MAX_TOKENS=1024
```

## Build

```sh
go build -o plugin-openai ./cmd/plugin
```

## Test

```sh
go test ./...
```

## Capability

`ai.complete`

Input:

```json
{
  "prompt": "Write a short reply",
  "model": "optional-model",
  "temperature": 0.7,
  "max_tokens": 1024
}
```

Output:

```json
{
  "text": "Generated response",
  "model": "model-name"
}
```

## Security

Never commit `.env` or API keys.
