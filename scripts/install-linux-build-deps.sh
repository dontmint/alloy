#!/usr/bin/env bash
# Install native build dependencies inside the CI compat container (AlmaLinux 8).
set -euo pipefail

if [[ -f /etc/os-release ]]; then
  # shellcheck disable=SC1091
  source /etc/os-release
else
  echo "Cannot detect OS; /etc/os-release missing" >&2
  exit 1
fi

case "${ID:-}" in
  almalinux|rhel|centos|rocky)
    dnf install -y \
      gcc gcc-c++ make git file xz \
      binutils \
      systemd-devel
    # Full profile UI build (Node 20 module stream on EL8).
    if [[ "${INSTALL_NODE:-0}" == "1" ]]; then
      dnf module install -y nodejs:20 || dnf install -y nodejs npm
    fi
    ;;
  debian|ubuntu)
    export DEBIAN_FRONTEND=noninteractive
    apt-get update
    apt-get install -y \
      build-essential make git file xz-utils \
      binutils \
      libsystemd-dev
    if [[ "${INSTALL_NODE:-0}" == "1" ]]; then
      apt-get install -y nodejs npm
    fi
    rm -rf /var/lib/apt/lists/*
    ;;
  *)
    echo "Unsupported build container OS: ${ID:-unknown}" >&2
    exit 1
    ;;
esac
