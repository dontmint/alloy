#!/usr/bin/env bash
# Build linux/amd64 Alloy inside an AlmaLinux 8 container (glibc 2.28 baseline).
#
# Usage:
#   ./scripts/build-linux-amd64-container.sh [full|slim]
#
# Requires Docker on the host. Matches the GitHub Actions release build environment.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROFILE="${1:-full}"

BUILD_IMAGE="${ALLOY_BUILD_CONTAINER_IMAGE:-almalinux:8}"
BUILD_SCRIPT="./scripts/build-linux-amd64.sh"
INSTALL_NODE=0
if [[ "$PROFILE" == "slim" ]]; then
  BUILD_SCRIPT="./scripts/build-linux-amd64-slim.sh"
else
  INSTALL_NODE=1
fi

GO_VERSION="$(grep '^go ' "${ROOT}/go.mod" | awk '{print $2}')"

docker run --rm \
  -v "${ROOT}:/src" \
  -w /src \
  -e RELEASE_BUILD="${RELEASE_BUILD:-1}" \
  -e VERSION="${VERSION:-dev}" \
  -e GOOS=linux \
  -e GOARCH=amd64 \
  -e CGO_ENABLED="${CGO_ENABLED:-1}" \
  -e INSTALL_NODE="${INSTALL_NODE}" \
  -e VERIFY_GLIBC="${VERIFY_GLIBC:-1}" \
  -e GLIBC_BASELINE_MAX="${GLIBC_BASELINE_MAX:-2.28}" \
  "${BUILD_IMAGE}" \
  bash -lc "
    ./scripts/install-linux-build-deps.sh
    curl -fsSL \"https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz\" \
      | tar -C /usr/local -xzf -
    export PATH=/usr/local/go/bin:\$PATH
    ${BUILD_SCRIPT}
  "
