# AGENTS.md

This repository is the official AuxiTalk OpenAI-compatible AI plugin.

It exposes AI capabilities to AuxiTalk Core through JSON-RPC 2.0 over line-delimited stdio.

## Required context

Read first:

1. `README.md`
2. `docs/ai-development-guide.md`
3. `plugin.json`
4. `internal/plugin/*`
5. `internal/openai/*`

## Required checks

Before finishing code changes:

```sh
gofmt -w <changed-go-files>
go test ./...
```

## Protocol rules

- stdout is reserved for JSON-RPC only.
- stderr is for logs.
- implement `plugin.health`.
- handle `plugin.stop` if applicable.
- never log API keys or request secrets.

## Safety rules

- Do not commit `.env` files or API keys.
- Avoid logging prompts if they may contain sensitive data.
- Keep OpenAI-compatible provider behavior configurable by env.
- Do not make network calls in unit tests unless explicitly designed as integration tests.

## Product framing

This plugin is one possible AI provider. AuxiTalk must remain provider-agnostic and usable with local or OpenAI-compatible models.
