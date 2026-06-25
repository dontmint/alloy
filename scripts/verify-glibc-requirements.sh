#!/usr/bin/env bash
# Print the highest GLIBC symbol version referenced by a Linux ELF binary.
set -euo pipefail

BIN="${1:-build/alloy}"

if [[ ! -f "$BIN" ]]; then
  echo "Binary not found: $BIN" >&2
  exit 1
fi

if ! command -v objdump >/dev/null 2>&1; then
  echo "objdump not available; skipping GLIBC verification"
  exit 0
fi

max="$(
  objdump -T "$BIN" 2>/dev/null \
    | grep -o 'GLIBC_[0-9.]*' \
    | sed 's/^GLIBC_//' \
    | sort -u -t. -k1,1n -k2,2n -k3,3n \
    | tail -1
)"

if [[ -z "$max" ]]; then
  echo "No dynamic GLIBC symbols found in $BIN (static binary or objdump unavailable)"
  exit 0
fi

count="$(
  objdump -T "$BIN" 2>/dev/null \
    | grep -o 'GLIBC_[0-9.]*' \
    | sed 's/^GLIBC_//' \
    | sort -u \
    | wc -l \
    | tr -d ' '
)"

echo "GLIBC baseline check for ${BIN}:"
printf '  required max GLIBC_%s\n' "$max"
echo "  (${count} distinct GLIBC symbol versions referenced)"

# Fail CI when a binary needs a newer glibc than our documented baseline (2.28).
baseline="${GLIBC_BASELINE_MAX:-2.28}"
highest="$(printf '%s\n' "$baseline" "$max" | sort -t. -k1,1n -k2,2n -k3,3n | tail -1)"
if [[ "$highest" != "$baseline" ]]; then
  echo "ERROR: GLIBC_${max} exceeds documented baseline GLIBC_${baseline}" >&2
  exit 1
fi
