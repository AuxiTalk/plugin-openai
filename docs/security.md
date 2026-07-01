# Security

## Secrets

Never commit API keys or `.env` files.

Use local environment variables or a secret manager.

## Logging

The plugin logs to stderr and must never log `OPENAI_API_KEY`.

## Network

This plugin sends prompts to the configured OpenAI-compatible API.

Use local providers such as Ollama or LM Studio if the prompt must remain on your machine.

## Timeouts

All API calls use a timeout configured by `OPENAI_TIMEOUT`.

Default: `60s`.
