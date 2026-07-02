# AI Development Guide

This guide is for AI coding agents working on the AuxiTalk OpenAI-compatible plugin.

## Responsibilities

- Provide AI completion capabilities to AuxiTalk Core.
- Remain compatible with OpenAI-style `/chat/completions` APIs.
- Keep provider configuration in env vars.
- Avoid leaking prompts, responses, or API keys in logs.

## Safe workflow

1. Inspect `plugin.json` and `internal/plugin`.
2. Add tests for runtime/capability behavior.
3. Mock HTTP where possible.
4. Run `gofmt` and `go test ./...`.
5. Commit only when requested.
6. Push only when explicitly requested.

## Sensitive areas

- HTTP client retries/timeouts.
- API key handling.
- Prompt/message logging.
- Provider compatibility.
