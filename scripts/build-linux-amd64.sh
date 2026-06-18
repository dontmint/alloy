#!/usr/bin/env bash
# Build linux/amd64 Alloy with embedded UI (no SKIP_UI_BUILD).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

export RELEASE_BUILD="${RELEASE_BUILD:-1}"
export GOOS="${GOOS:-linux}"
export GOARCH="${GOARCH:-amd64}"
export GO_TAGS="${GO_TAGS:-gore2regex embedalloyui}"

echo "Building Alloy ${GOOS}/${GOARCH} with GO_TAGS=${GO_TAGS}"
make alloy

ls -lh build/alloy
file build/alloy
