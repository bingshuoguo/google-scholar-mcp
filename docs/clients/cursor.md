# Cursor

This guide configures `google-scholar-mcp` as a local MCP server in Cursor.

## Prerequisites

- Cursor installed
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

Cursor supports MCP config in either of these locations:

- `~/.cursor/mcp.json` for all projects
- `.cursor/mcp.json` inside a specific workspace

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

Use the absolute binary path if Cursor does not inherit your shell `PATH`.

## Verify In Cursor

1. Open Cursor settings and the MCP integrations page.
2. Confirm `google-scholar` appears in the server list.
3. Refresh if needed.
4. Ask Cursor to call `search_google_scholar_key_words` with a simple query.

## Troubleshooting

- JSON syntax errors will prevent Cursor from loading any MCP servers in that file.
- If the server does not start, replace `command` with the absolute binary path.
- Keep logs on `stderr`. Do not wrap the binary with a script that prints to `stdout`.

## Reference

- Cursor MCP docs: https://docs.cursor.com/en/context/mcp
