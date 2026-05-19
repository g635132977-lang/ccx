# Connect Claude Code to CCX

Claude Code uses the Anthropic Messages protocol. When connecting it to CCX, use the CCX **Messages** endpoint.

## How it works

```text
Claude Code  ->  CCX /v1/messages  ->  Messages channel  ->  upstream Anthropic-compatible endpoint
```

## 1. Configure a CCX channel

1. Open the CCX admin console and go to **Messages**
2. Click "Add Channel"
3. Fill in the channel configuration for your upstream provider

Common examples:

| Upstream | Service type | Base URL example |
|----------|--------------|------------------|
| Anthropic Claude | `Claude` | `https://api.anthropic.com` |
| DeepSeek Anthropic-compatible | `Claude` | `https://api.deepseek.com/anthropic` |
| Kimi coding endpoint | `Claude` | `https://api.kimi.com/coding/` |
| GLM Anthropic-compatible | `Claude` | `https://open.bigmodel.cn/api/anthropic` |

::: tip
For upstream API keys, model names, and special options, see the matching [provider setup guide](/en/providers/).
:::

## 2. Configure Claude Code

Set these environment variables:

```bash
export ANTHROPIC_API_KEY="your-ccx-proxy-key"
export ANTHROPIC_BASE_URL="http://localhost:3000"
```

Then run:

```bash
claude "hello"
```

::: warning
`ANTHROPIC_BASE_URL` should be the CCX gateway root. Do not append `/v1` or `/v1/messages`.
:::

## 3. Recommended model mapping

Claude Code may send Claude model names such as:

- `claude-opus-4-7`
- `claude-sonnet-4-6`
- `claude-haiku-4-5-20251001`

If the upstream is not Anthropic Claude, configure model mapping on the Messages channel:

| Request model match | Example upstream model |
|---------------------|------------------------|
| `opus` | upstream high-capability model |
| `sonnet` | upstream primary model |
| `haiku` | upstream lightweight model |

CCX prefers longer matching keys first. Use the actual upstream model names from your provider.

## Troubleshooting

### Claude Code returns 401 Unauthorized

Check:

1. `ANTHROPIC_API_KEY` matches CCX `PROXY_ACCESS_KEY`
2. CCX was started with the expected `PROXY_ACCESS_KEY`
3. No old shell environment variable overrides the current value

### Claude Code still calls official Claude

Check the current authentication and environment. The terminal session running Claude Code must see both variables:

```bash
printenv ANTHROPIC_BASE_URL
printenv ANTHROPIC_API_KEY
```

### 404 or missing endpoint

The usual cause is an incorrect `ANTHROPIC_BASE_URL`.

Correct:

```bash
export ANTHROPIC_BASE_URL="http://localhost:3000"
```

Incorrect examples:

```bash
export ANTHROPIC_BASE_URL="http://localhost:3000/v1"
export ANTHROPIC_BASE_URL="http://localhost:3000/v1/messages"
```

### model_not_found

Check the Messages channel:

1. The model allowlist includes the requested model or mapped upstream model
2. Model mapping covers the actual model name requested by Claude Code
3. The upstream model name is correct

### thinking or signature-related errors

If the upstream is not fully compatible with Claude thinking/signature behavior:

1. Prefer an upstream with real Claude Messages compatibility
2. Check reasoning/thinking compatibility options on the channel
3. Avoid routing long thinking-enabled sessions across incompatible upstreams

### Claude Code `/model` does not show expected models

Confirm CCX has at least one available Messages channel and the model list includes the expected models. You can also verify the CCX models endpoint:

```bash
curl http://localhost:3000/v1/models \
  -H "Authorization: Bearer your-ccx-proxy-key"
```
