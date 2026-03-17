# Installation

Recommended order:

1. `go install` if you already use Go
2. GitHub Release archives if you want a direct download
3. Homebrew or Scoop if you prefer package-manager-based installation

## Go Install

This is the most direct installation path for Go users:

```bash
go install github.com/bingshuoguo/google-scholar-mcp/cmd/google-scholar-mcp@latest
google-scholar-mcp --version
```

If you want a fixed released version instead of `latest`, pin the tag explicitly:

```bash
go install github.com/bingshuoguo/google-scholar-mcp/cmd/google-scholar-mcp@v0.1.1
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

After extraction:

- macOS and Linux: make sure `google-scholar-mcp` is executable and on your `PATH`
- Windows: make sure `google-scholar-mcp.exe` is on your `PATH`

## Homebrew

Install from the formula tracked in this repository:

```bash
brew install --formula https://raw.githubusercontent.com/bingshuoguo/google-scholar-mcp/main/packaging/homebrew/google-scholar-mcp.rb
```

Verify:

```bash
google-scholar-mcp --version
```

Typical Homebrew-installed paths include:

- Apple Silicon macOS: `/opt/homebrew/bin/google-scholar-mcp`
- Intel macOS: `/usr/local/bin/google-scholar-mcp`

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

Useful checks:

```bash
command -v google-scholar-mcp
google-scholar-mcp help
google-scholar-mcp --version
```
