# Claude

This guide configures `google-scholar-mcp` for Claude clients.

## Prerequisites

- Claude Code or Claude Desktop installed
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

## Claude Code CLI

Anthropic documents the `claude mcp add` workflow for local stdio servers.

Example:

```bash
claude mcp add google-scholar /absolute/path/to/google-scholar-mcp/.bin/google-scholar-mcp
```

## Claude Desktop Manual Configuration

If you are using Claude Desktop, edit `claude_desktop_config.json` and add:

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

Common config locations:

- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\\Claude\\claude_desktop_config.json`

## Verify In Claude

- In Claude Code, run `claude mcp list`.
- In Claude Desktop, restart the app after saving the config.
- Ask Claude to call `search_google_scholar_key_words`.

## Troubleshooting

- Use the absolute path to the binary if the client cannot resolve it.
- Do not add wrappers that print to `stdout`.
- If a config change is ignored, restart Claude completely.

## Reference

- Anthropic Claude Code MCP docs: https://docs.anthropic.com/en/docs/claude-code/mcp
