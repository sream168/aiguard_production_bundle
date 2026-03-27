#!/usr/bin/env bash
set -euo pipefail

resolve_wails_bin() {
  if command -v wails >/dev/null 2>&1; then
    command -v wails
    return 0
  fi

  if command -v go >/dev/null 2>&1; then
    local gopath
    gopath="$(go env GOPATH 2>/dev/null || true)"
    if [ -n "${gopath:-}" ] && [ -x "$gopath/bin/wails" ]; then
      printf '%s\n' "$gopath/bin/wails"
      return 0
    fi
  fi

  return 1
}

resolve_linux_wails_tags() {
  if [ "$(uname -s)" != "Linux" ]; then
    return 0
  fi

  if command -v pkg-config >/dev/null 2>&1 && pkg-config --exists webkit2gtk-4.1; then
    printf '%s\n' "webkit2_41"
  fi
}
