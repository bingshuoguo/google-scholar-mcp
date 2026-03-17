# Codex

This guide configures `google-scholar-mcp` for OpenAI Codex.

## Prerequisites

- Codex installed
- Go `1.25+`
- This repository available locally

Build the server binary:

```bash
go build -o ./.bin/google-scholar-mcp ./cmd/google-scholar-mcp
```

Optional validation:

```bash
./scripts/verify_stdio.sh smoke
```

## Configuration File

Codex uses TOML, not JSON.

Config file:

- `~/.codex/config.toml`

## Example Configuration

```toml
[mcp_servers.google-scholar]
command = "/absolute/path/to/google-scholar-mcp/.bin/google-scholar-mcp"
args = []
env = { LOG_LEVEL = "info", SCHOLAR_MAX_RESULTS = "5" }
```

## Verify In Codex

```bash
codex mcp list
codex mcp get google-scholar
```

Then start Codex and ask it to search Scholar.

## Troubleshooting

- The key is `mcp_servers`, not `mcp-servers`.
- TOML syntax errors in `~/.codex/config.toml` can break all Codex MCP entries.
- If Codex cannot find the binary, switch `command` to an absolute path.

## Reference

- OpenAI Codex MCP docs: https://developers.openai.com/codex/cli#model-context-protocol-mcp
