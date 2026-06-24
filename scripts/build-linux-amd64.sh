#!/usr/bin/env bash
# Build linux/amd64 Alloy (full profile by default; set ALLOY_BUILD_PROFILE=slim for slim).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

export ALLOY_BUILD_PROFILE="${ALLOY_BUILD_PROFILE:-full}"
export RELEASE_BUILD="${RELEASE_BUILD:-1}"
export GOOS="${GOOS:-linux}"
export GOARCH="${GOARCH:-amd64}"

if [[ "${ALLOY_BUILD_PROFILE}" == "slim" ]]; then
  export SKIP_UI_BUILD="${SKIP_UI_BUILD:-1}"
  export GO_TAGS="${GO_TAGS:-gore2regex alloy_slim}"
else
  export GO_TAGS="${GO_TAGS:-gore2regex embedalloyui}"
fi

echo "Building Alloy (${ALLOY_BUILD_PROFILE}) ${GOOS}/${GOARCH} with GO_TAGS=${GO_TAGS}"
make alloy

ls -lh build/alloy
file build/alloy
