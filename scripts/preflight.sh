#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
source "$ROOT/scripts/lib_wails.sh"

check_cmd() {
  if command -v "$1" >/dev/null 2>&1; then
    echo "[ OK ] $1 -> $(command -v "$1")"
  else
    echo "[FAIL] missing command: $1" >&2
    return 1
  fi
}

check_cmd go
check_cmd git
check_cmd npm

if WAILS_BIN="$(resolve_wails_bin)"; then
  echo "[ OK ] wails -> $WAILS_BIN"
else
  echo "[FAIL] missing command: wails" >&2
  exit 1
fi

if [ "$(uname -s)" = "Linux" ] && command -v pkg-config >/dev/null 2>&1; then
  if pkg-config --exists webkit2gtk-4.1; then
    echo "[ OK ] detected webkit2gtk-4.1 via pkg-config"
    echo "[INFO] Ubuntu 24.04+ should build with Wails tag: webkit2_41"
  elif pkg-config --exists webkit2gtk-4.0; then
    echo "[ OK ] detected webkit2gtk-4.0 via pkg-config"
  else
    echo "[WARN] neither webkit2gtk-4.0 nor webkit2gtk-4.1 was found via pkg-config" >&2
  fi
fi

if command -v xcode-select >/dev/null 2>&1; then
  if xcode-select -p >/dev/null 2>&1; then
    echo "[ OK ] Xcode Command Line Tools detected"
  else
    echo "[WARN] Xcode Command Line Tools are not installed. Run: xcode-select --install" >&2
  fi
fi

echo
echo "running wails doctor..."
"$WAILS_BIN" doctor
