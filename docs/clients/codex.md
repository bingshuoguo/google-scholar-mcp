# Codex

This guide configures `google-scholar-mcp` for OpenAI Codex.

## Prerequisites

- Codex installed
- `google-scholar-mcp` installed locally

## Install the Binary

Recommended:

```bash
go install github.com/bingshuoguo/google-scholar-mcp/cmd/google-scholar-mcp@latest
google-scholar-mcp --version
```

Alternative installation paths are documented in [Installation](../install.md).

Optional local validation:

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
command = "google-scholar-mcp"
args = []
env = { LOG_LEVEL = "info", SCHOLAR_MAX_RESULTS = "5" }
```

If Codex does not resolve the binary from `PATH`, switch `command` to an absolute path such as:

```toml
[mcp_servers.google-scholar]
command = "/Users/your-name/go/bin/google-scholar-mcp"
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
- Confirm installation with `google-scholar-mcp --version` before debugging MCP config.

## Reference

- OpenAI Codex MCP docs: https://developers.openai.com/codex/cli#model-context-protocol-mcp
