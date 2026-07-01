# Configuration

## Environment variables

```txt
OPENAI_API_KEY                required
OPENAI_BASE_URL               optional, default https://api.openai.com/v1
OPENAI_MODEL                  optional, default gpt-4o-mini
OPENAI_MODEL_FAST             optional profile model
OPENAI_MODEL_SMART            optional profile model
OPENAI_MODEL_LOCAL            optional profile model
OPENAI_TIMEOUT                optional, default 60s
OPENAI_TEMPERATURE            optional
OPENAI_MAX_TOKENS             optional
OPENAI_TOP_P                  optional
OPENAI_PRESENCE_PENALTY       optional
OPENAI_FREQUENCY_PENALTY      optional
OPENAI_RETRIES                optional, default 0
OPENAI_RETRY_BACKOFF          optional, default 1s
OPENAI_HEALTH_CHECK           optional, default false
```

All advanced values are optional. By default, retries and real health checks are disabled.

## Model profiles

The input can pass `profile` instead of a direct model:

```json
{
  "profile": "fast",
  "prompt": "Reply briefly"
}
```

Supported profile names:

```txt
fast   -> OPENAI_MODEL_FAST
smart  -> OPENAI_MODEL_SMART
local  -> OPENAI_MODEL_LOCAL
```

If a profile is missing, the plugin falls back to `OPENAI_MODEL`.

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
