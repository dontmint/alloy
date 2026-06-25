#!/usr/bin/env bash
# Build linux/amd64 Alloy using the slim build profile.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

export ALLOY_BUILD_PROFILE=slim
export RELEASE_BUILD="${RELEASE_BUILD:-1}"
export GOOS="${GOOS:-linux}"
export GOARCH="${GOARCH:-amd64}"
export SKIP_UI_BUILD="${SKIP_UI_BUILD:-1}"
# Slim profile omits embedalloyui; allow extra tags via GO_TAGS.
export GO_TAGS="${GO_TAGS:-gore2regex alloy_slim}"

echo "Building slim Alloy ${GOOS}/${GOARCH} with GO_TAGS=${GO_TAGS} CGO_ENABLED=${CGO_ENABLED:-1}"
make alloy

ls -lh build/alloy
file build/alloy

if [[ "${VERIFY_GLIBC:-0}" == "1" ]]; then
  ./scripts/verify-glibc-requirements.sh build/alloy
fi
