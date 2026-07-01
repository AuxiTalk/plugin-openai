# Configuration

## Environment variables

```txt
OPENAI_API_KEY       required
OPENAI_BASE_URL      optional, default https://api.openai.com/v1
OPENAI_MODEL         optional, default gpt-4o-mini
OPENAI_TIMEOUT       optional, default 60s
OPENAI_TEMPERATURE   optional
OPENAI_MAX_TOKENS    optional
```

## Provider examples

### OpenAI

```txt
OPENAI_API_KEY=sk-...
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-4o-mini
```

### Ollama

```txt
OPENAI_API_KEY=ollama
OPENAI_BASE_URL=http://localhost:11434/v1
OPENAI_MODEL=llama3.1
```

### LM Studio

```txt
OPENAI_API_KEY=lm-studio
OPENAI_BASE_URL=http://localhost:1234/v1
OPENAI_MODEL=local-model
```
