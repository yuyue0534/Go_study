#!/usr/bin/env sh
set -eu
cd "$(dirname "$0")/.."
echo "[1/2] Downloading dependencies..."
go mod tidy
echo "[2/2] Starting student system..."
go run ./cmd/server
