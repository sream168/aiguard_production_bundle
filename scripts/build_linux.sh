#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

source "$ROOT/scripts/lib_wails.sh"

bash ./scripts/preflight.sh

WAILS_BIN="$(resolve_wails_bin)"
LINUX_TAGS="$(resolve_linux_wails_tags || true)"

echo
echo "[1/5] install frontend deps"
cd frontend
npm install
npm run build
cd ..

echo
echo "[2/5] tidy go modules"
go mod tidy

echo
echo "[3/5] run wails doctor"
"$WAILS_BIN" doctor || true

echo
echo "[4/5] build Linux binary"
if [ -n "${LINUX_TAGS:-}" ]; then
  echo "detected webkit2gtk-4.1; use build tags: $LINUX_TAGS"
  "$WAILS_BIN" build -platform linux/amd64 -tags "$LINUX_TAGS"
else
  "$WAILS_BIN" build -platform linux/amd64
fi

echo
echo "[5/5] package release"
bash ./scripts/package_release_linux.sh

echo
echo "done. artifacts are in build/bin and release"
