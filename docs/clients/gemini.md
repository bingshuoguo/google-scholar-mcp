# Gemini CLI

This guide configures `google-scholar-mcp` for Gemini CLI.

## Prerequisites

- Gemini CLI installed
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

Gemini CLI reads user settings from:

- `~/.gemini/settings.json`

## Example Configuration

```json
{
  "mcpServers": {
    "google-scholar": {
      "command": "/absolute/path/to/google-scholar-mcp/.bin/google-scholar-mcp",
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

## Reference

- Gemini CLI MCP docs: https://github.com/google-gemini/gemini-cli/blob/main/docs/tools/mcp-server.md
