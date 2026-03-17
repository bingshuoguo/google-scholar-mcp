# Installation

## Go Install

This is the most direct installation path for Go users:

```bash
go install github.com/bingshuoguo/google-scholar-mcp/cmd/google-scholar-mcp@latest
google-scholar-mcp --version
```

## GitHub Releases

Release archives are published for:

- macOS `amd64`
- macOS `arm64`
- Linux `amd64`
- Linux `arm64`
- Windows `amd64`
- Windows `arm64`

Download the matching archive from the GitHub Releases page, extract it, and put `google-scholar-mcp` or `google-scholar-mcp.exe` in your `PATH`.

Current release page:

- https://github.com/bingshuoguo/google-scholar-mcp/releases/tag/v0.1.1

## Homebrew

Install from the formula tracked in this repository:

```bash
brew install --formula https://raw.githubusercontent.com/bingshuoguo/google-scholar-mcp/main/packaging/homebrew/google-scholar-mcp.rb
```

Verify:

```bash
google-scholar-mcp --version
```

## Scoop

Install from the Scoop manifest tracked in this repository:

```powershell
scoop install https://raw.githubusercontent.com/bingshuoguo/google-scholar-mcp/main/packaging/scoop/google-scholar-mcp.json
```

Verify:

```powershell
google-scholar-mcp --version
```

## After Installation

Once the binary is installed, register it with your MCP client. For example:

- Codex: use `command = "google-scholar-mcp"` if the binary is in `PATH`
- Cursor, Claude Desktop, Gemini CLI: prefer an absolute binary path if GUI apps do not inherit your shell `PATH`
