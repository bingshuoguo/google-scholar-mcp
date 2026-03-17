#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$ROOT_DIR/.bin"
SERVER_BIN="$BIN_DIR/google-scholar-mcp"
MODE="${1:-smoke}"
TEMP_CONFIG=""

cleanup() {
  if [[ -n "${TEMP_CONFIG:-}" && -f "${TEMP_CONFIG:-}" ]]; then
    rm -f "$TEMP_CONFIG"
  fi
}

trap cleanup EXIT

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

json_escape() {
  local value="$1"
  value="${value//\\/\\\\}"
  value="${value//\"/\\\"}"
  value="${value//$'\n'/\\n}"
  value="${value//$'\r'/\\r}"
  value="${value//$'\t'/\\t}"
  printf '%s' "$value"
}

build_server() {
  require_cmd go
  mkdir -p "$BIN_DIR"
  echo "Building Go server at $SERVER_BIN"
  (
    cd "$ROOT_DIR"
    go build -o "$SERVER_BIN" ./cmd/google-scholar-mcp
  )
}

make_config() {
  local config_file
  config_file="$(mktemp)"
  cat >"$config_file" <<EOF
{
  "mcpServers": {
    "default-server": {
      "command": "$(json_escape "$SERVER_BIN")",
      "args": [],
      "env": {
        "LOG_LEVEL": "error"
      }
    }
  }
}
EOF
  printf '%s\n' "$config_file"
}

run_smoke() {
  require_cmd go
  echo "Running stdio smoke test with the MCP Go SDK client"
  (
    cd "$ROOT_DIR"
    go run ./scripts/smoke_stdio "$SERVER_BIN"
  )
}

run_ui() {
  require_cmd npx
  TEMP_CONFIG="$(make_config)"

  echo "Launching MCP Inspector UI for $SERVER_BIN"
  echo "Inspector UI default URL: http://localhost:6274"
  npx -y @modelcontextprotocol/inspector --config "$TEMP_CONFIG"
}

build_server

case "$MODE" in
  smoke)
    run_smoke
    ;;
  ui)
    run_ui
    ;;
  *)
    echo "Usage: $0 [smoke|ui]" >&2
    exit 1
    ;;
esac
