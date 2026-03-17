#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
  echo "Usage: $0 <version>" >&2
  echo "Example: $0 v0.1.1" >&2
  exit 1
fi

VERSION="${VERSION#v}"
TAG="v${VERSION}"
DIST_DIR="$ROOT_DIR/dist/$TAG"

mkdir -p "$DIST_DIR"
rm -f "$DIST_DIR"/*

build_archive() {
  local goos="$1"
  local goarch="$2"
  local ext="$3"
  local binary_name="google-scholar-mcp"
  local archive_name="google-scholar-mcp_${VERSION}_${goos}_${goarch}.${ext}"
  local work_dir
  work_dir="$(mktemp -d)"

  if [[ "$goos" == "windows" ]]; then
    binary_name="${binary_name}.exe"
  fi

  echo "Building ${goos}/${goarch}"
  (
    cd "$ROOT_DIR"
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
      go build -trimpath -ldflags="-s -w -X main.version=${TAG}" \
      -o "${work_dir}/${binary_name}" ./cmd/google-scholar-mcp
  )

  if [[ "$ext" == "zip" ]]; then
    (
      cd "$work_dir"
      zip -q "$DIST_DIR/$archive_name" "$binary_name"
    )
  else
    tar -C "$work_dir" -czf "$DIST_DIR/$archive_name" "$binary_name"
  fi

  rm -rf "$work_dir"
}

build_archive darwin amd64 tar.gz
build_archive darwin arm64 tar.gz
build_archive linux amd64 tar.gz
build_archive linux arm64 tar.gz
build_archive windows amd64 zip
build_archive windows arm64 zip

(
  cd "$DIST_DIR"
  shasum -a 256 * > checksums.txt
)

echo "Release artifacts written to $DIST_DIR"
