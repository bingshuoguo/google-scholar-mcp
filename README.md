# Google Scholar MCP

`google-scholar-mcp` is a local `stdio` MCP server written in Go. It retrieves Google Scholar search results directly from Scholar HTML and exposes them as structured tools for MCP clients such as Cursor, Codex, Claude, and Gemini CLI.

This repository is designed for low-volume local usage. It does not depend on a paid SERP API, but it also does not try to hide the fact that Google Scholar scraping is brittle and best-effort.

## Features

- Go implementation built on the official MCP Go SDK
- Local `stdio` transport for desktop and CLI MCP clients
- Structured Google Scholar search results
- Best-effort Google Scholar author profile lookup
- Fixture-based parser tests
- Local smoke test and MCP Inspector validation script

## Available Tools

| Tool | Purpose |
| --- | --- |
| `search_google_scholar_key_words` | Search Google Scholar by keyword query and return structured paper metadata. |
| `search_google_scholar_advanced` | Search Google Scholar with author and year filters. |
| `get_author_info` | Find a Google Scholar author profile by name and return structured author metadata. |

## What Data You Can Get

The server can usually extract these fields when they are present on Google Scholar pages:

- paper title
- result link
- author and venue line
- Scholar result snippet
- year when detectable
- citation count
- version count
- PDF or resource link when present
- author profile metadata and publication summaries

Important boundary:

- The Scholar snippet is a search-result excerpt, not a guaranteed full abstract.
- Full paper abstracts, PDFs, and publisher-only metadata are out of scope unless they already appear on the Scholar page being scraped.

## Quick Start

### Requirements

- Go `1.23+`
- Node.js only if you want to launch MCP Inspector via `npx`

### Install

Install the binary from GitHub:

```bash
go install github.com/bingshuoguo/google-scholar-mcp/cmd/google-scholar-mcp@latest
```

Or build from source:

```bash
git clone git@github.com:bingshuoguo/google-scholar-mcp.git
cd google-scholar-mcp
go build -o ./.bin/google-scholar-mcp ./cmd/google-scholar-mcp
```

### Run

```bash
./.bin/google-scholar-mcp
```

For local development:

```bash
go run ./cmd/google-scholar-mcp
```

## Client Integration

- [Cursor](docs/clients/cursor.md)
- [Codex](docs/clients/codex.md)
- [Claude](docs/clients/claude.md)
- [Gemini CLI](docs/clients/gemini.md)

## Local Validation

Run the Go unit tests:

```bash
GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn go test ./...
```

Run the local `stdio` smoke test:

```bash
./scripts/verify_stdio.sh smoke
```

Launch MCP Inspector:

```bash
./scripts/verify_stdio.sh ui
```

The Inspector UI defaults to `http://localhost:6274`.

## Configuration

The server is configured with environment variables.

| Variable | Default | Purpose |
| --- | --- | --- |
| `MCP_TRANSPORT` | `stdio` | MCP transport. V1 is intended for local `stdio` use. |
| `SCHOLAR_BASE_URL` | `https://scholar.google.com` | Base Google Scholar URL. |
| `SCHOLAR_TIMEOUT` | `15s` | Upstream HTTP timeout. |
| `SCHOLAR_MAX_RESULTS` | `10` | Default maximum result count. |
| `SCHOLAR_RATE_LIMIT_RPS` | `0.5` | Request rate limit for upstream Scholar requests. |
| `SCHOLAR_USER_AGENT` | built-in default | HTTP user agent for Scholar requests. |
| `SCHOLAR_ACCEPT_LANGUAGE` | built-in default | Accept-Language header for Scholar requests. |
| `SCHOLAR_ENABLE_AUTHOR_TOOL` | `true` | Enable or disable `get_author_info`. |
| `LOG_LEVEL` | `info` | Structured logging level. |

Example:

```bash
LOG_LEVEL=debug SCHOLAR_MAX_RESULTS=5 ./.bin/google-scholar-mcp
```

## Development

### Project Layout

- `cmd/google-scholar-mcp`: executable entrypoint
- `internal/config`: configuration and logger setup
- `internal/mcpserver`: MCP server wiring and tool registration
- `internal/model`: shared domain models
- `internal/scholar`: Google Scholar provider, HTTP client, parsers, and tests
- `testdata`: HTML fixtures for parser tests
- `scripts/verify_stdio.sh`: build, smoke test, and Inspector launcher
- `scripts/smoke_stdio`: Go-based local smoke test helper

### Notes

- Logs are written to `stderr`, not `stdout`, so the server can safely run over `stdio`.
- Tool names are kept compatible with the legacy Python implementation.
- This repository currently focuses on local interactive usage, not bulk harvesting.

## Docs

- [Go rewrite design notes](docs/design.md)
- [Cursor integration](docs/clients/cursor.md)
- [Codex integration](docs/clients/codex.md)
- [Claude integration](docs/clients/claude.md)
- [Gemini CLI integration](docs/clients/gemini.md)

## Limitations

- Google Scholar has no stable public API for this workflow.
- HTML scraping can break when Scholar changes its markup.
- High-volume crawling, anti-bot challenge handling, and full-text retrieval are intentionally out of scope.

## Responsible Use

Use this project carefully and at low volume. You are responsible for respecting Google Scholar's terms of service, robots behavior, and any applicable usage limits in your environment.
