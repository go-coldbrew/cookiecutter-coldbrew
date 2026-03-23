#!/usr/bin/env bash
set -euxo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
TMPDIR="$(mktemp -d)"

cleanup() {
    echo "Cleaning up $TMPDIR"
    rm -rf "$TMPDIR"
}
trap cleanup EXIT

echo "=== Generating project in $TMPDIR ==="
cookiecutter --no-input --verbose --output-dir "$TMPDIR" "$REPO_DIR"

PROJECT_DIR="$TMPDIR/myapp"
cd "$PROJECT_DIR"

echo "=== make build ==="
make build

echo "=== make test ==="
make test

echo "=== make lint ==="
make lint

echo "=== All checks passed ==="
