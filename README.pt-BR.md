# AuxiTalk OpenAI Plugin

Plugin de IA compatível com OpenAI para o AuxiTalk.

> English documentation: [README.md](README.md)

## Visão geral

Este plugin expõe a capability `ai.complete` e conecta o AuxiTalk a APIs de chat completions compatíveis com OpenAI.

Provedores compatíveis incluem:

- OpenAI
- Groq
- Together.ai
- Ollama
- LM Studio
- qualquer API compatível com `/chat/completions` da OpenAI

## Configuração

Use variáveis de ambiente:

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

## Testes

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

## Segurança

Nunca commite `.env` ou chaves de API.
