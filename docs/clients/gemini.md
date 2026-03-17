# Gemini CLI

This guide configures `google-scholar-mcp` for Gemini CLI.

## Prerequisites

- Gemini CLI installed
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

Gemini CLI reads user settings from:

- `~/.gemini/settings.json`

## Example Configuration

```json
{
  "mcpServers": {
    "google-scholar": {
      "command": "/absolute/path/to/google-scholar-mcp",
      "args": [],
      "env": {
        "LOG_LEVEL": "info",
        "SCHOLAR_MAX_RESULTS": "5"
      }
    }
  }
}
```

## Verify In Gemini CLI

Start Gemini CLI and ask it to use `search_google_scholar_key_words` for a simple query.

If the server is not detected, re-open Gemini CLI after updating `settings.json`.

## Troubleshooting

- `settings.json` must remain valid JSON.
- Use an absolute `command` path if Gemini CLI does not inherit your shell `PATH`.
- Keep the server on `stdio`; do not wrap it with a command that writes extra output.
- Confirm installation with `google-scholar-mcp --version` before debugging Gemini config.

## Reference

- Gemini CLI MCP docs: https://github.com/google-gemini/gemini-cli/blob/main/docs/tools/mcp-server.md
