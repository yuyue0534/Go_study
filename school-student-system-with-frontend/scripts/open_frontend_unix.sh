#!/usr/bin/env sh
set -eu
ROOT_DIR="$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)"
INDEX_FILE="$ROOT_DIR/frontend/index.html"

if command -v xdg-open >/dev/null 2>&1; then
  xdg-open "$INDEX_FILE"
elif command -v open >/dev/null 2>&1; then
  open "$INDEX_FILE"
else
  printf '%s\n' "Open this file in your browser: $INDEX_FILE"
fi
