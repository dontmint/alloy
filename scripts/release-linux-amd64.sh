#!/usr/bin/env bash
# Build, package, and upload linux/amd64 release asset (replaces existing file on the tag).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

TAG="${1:-}"
if [[ -z "$TAG" ]]; then
  echo "Usage: $0 <tag>" >&2
  echo "Example: $0 v1.17.0-docker-state" >&2
  exit 1
fi

export VERSION="$TAG"

echo "==> Building ${TAG} locally"
"$ROOT/scripts/build-linux-amd64.sh"

mkdir -p dist
cp build/alloy dist/alloy-linux-amd64
chmod +x dist/alloy-linux-amd64
xz -9 -k -f dist/alloy-linux-amd64
ls -lh dist/alloy-linux-amd64.xz

if ! command -v gh >/dev/null 2>&1; then
  echo "gh CLI not found; binary ready at dist/alloy-linux-amd64.xz"
  exit 0
fi

echo "==> Uploading dist/alloy-linux-amd64.xz to GitHub release ${TAG}"
gh release upload "$TAG" dist/alloy-linux-amd64.xz --clobber

echo "Done: https://github.com/$(gh repo view --json nameWithOwner -q .nameWithOwner)/releases/tag/${TAG}"
