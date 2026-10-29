#!/usr/bin/env bash
set -euo pipefail
test -s README.md
go test ./...
git diff --check
