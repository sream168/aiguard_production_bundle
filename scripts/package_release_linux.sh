#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

mkdir -p release
STAMP="$(date +%Y%m%d-%H%M%S)"
BIN_PATH="build/bin/AIGuard"
ARCHIVE_PATH="release/AIGuard-linux-${STAMP}.tar.gz"

if [ ! -f "$BIN_PATH" ]; then
  echo "build/bin/AIGuard not found, please build first" >&2
  exit 1
fi

tar -czf "$ARCHIVE_PATH" -C build/bin AIGuard
echo "packaged: $ARCHIVE_PATH"
