# Architecture

This repository is organized as a small, layered Go MCP server.

## Entry Point

- `cmd/google-scholar-mcp`
- handles CLI commands: `stdio`, `version`, `help`
- resolves build version and starts the stdio server

## MCP Layer

- `internal/mcpserver`
- owns MCP server metadata
- registers `tools`, `resources`, and `prompts`
- translates domain errors into model-friendly MCP responses

## Domain Layer

- `internal/scholar`
- builds Scholar requests
- performs HTTP retrieval
- parses Scholar HTML
- classifies upstream and parsing failures

## Shared Models

- `internal/model`
- search and author response types returned through MCP tools

## Validation

- `internal/scholar/*_test.go` validates parser behavior with HTML fixtures
- `scripts/smoke_stdio` checks the exposed MCP surface over stdio
- `scripts/verify_stdio.sh` provides local smoke and Inspector verification
