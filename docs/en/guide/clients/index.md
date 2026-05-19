# Client Setup Overview

CCX client setup docs are organized by client. Choose the client you use, then follow the matching page.

## Choose a client

| Client | Recommended endpoint | CCX Base URL | Use case |
|--------|----------------------|--------------|----------|
| [Claude Code](./claude-code) | Messages | `http://localhost:3000` | Coding assistant using Claude Messages protocol |
| [Codex CLI / Codex App](./codex) | Responses | `http://localhost:3000/v1` | Codex tools using OpenAI Responses protocol |
| [OpenCode](./opencode) | Chat | `http://localhost:3000/v1` | Coding tools using OpenAI Chat-compatible protocol |

## Before connecting a client

1. Start CCX and confirm the gateway URL, for example `http://localhost:3000`
2. Set `PROXY_ACCESS_KEY`
3. Add at least one working channel for the target endpoint in the CCX admin console
4. Use `PROXY_ACCESS_KEY` as the client API key

## Endpoint mapping

```text
Claude Code           ->  /v1/messages          ->  Messages channels
Codex CLI / App       ->  /v1/responses         ->  Responses channels
OpenCode              ->  /v1/chat/completions  ->  Chat channels
```

::: tip
Base URL rules differ by client. Claude Code uses the gateway root. Codex and OpenCode usually use the OpenAI-compatible `/v1` base URL.
:::

## General troubleshooting

For client-specific issues, open the matching client page.

| Symptom | Check |
|---------|-------|
| `401 Unauthorized` | The client API key must match CCX `PROXY_ACCESS_KEY` |
| `Connection refused` | CCX is running and listening on the expected port |
| `model_not_found` | Channel model allowlist and model mapping cover the requested model |
| Request never reaches a channel | The client is using the expected protocol endpoint |
