# Linux binary compatibility

Release assets `alloy-linux-amd64.xz` and `alloy-linux-amd64-slim.xz` are built on
**AlmaLinux 8** (glibc **2.28**, same baseline as RHEL 8).

## Supported distros (single linux-amd64 asset)

One `linux-amd64` binary is sufficient for current fleet and typical enterprise Linux:

| Distro | Minimum version | glibc |
| --- | --- | --- |
| RHEL / AlmaLinux / Rocky | 8+ | 2.28+ |
| Ubuntu | 20.04+ | 2.31+ |
| Debian | 10 (buster)+ | 2.28+ |
| SLES | 15+ | 2.28+ |

The binary is **not** expected to run on glibc 2.17 systems (RHEL 7 / CentOS 7) or
Ubuntu 18.04 without a separate older-glibc build pipeline.

## Why not `ubuntu-latest` or `ubuntu-22.04` runners?

GitHub `ubuntu-latest` (24.04) links against **GLIBC_2.38+**. Even `ubuntu-22.04`
runners produce binaries requiring **GLIBC_2.35**, which matches jammy but is
unnecessarily new for RHEL and prevents a single artifact from covering EL8.

Building on AlmaLinux 8 links against glibc 2.28 while remaining compatible with
all newer glibc versions (glibc maintains backward compatibility).

## CGO and custom components

| Component | CGO required? | Notes |
| --- | --- | --- |
| `prometheus.exporter.docker_state` | No | Pure Go + Docker API |
| `prometheus.exporter.smartctl` | No | Shells out to `smartctl` binary on host |
| `otelcol.receiver.sqlquery` | No | `database/sql` + pure Go drivers (ClickHouse, Postgres) |
| `beyla.ebpf` | No (Go eBPF) | Requires **kernel** BTF/eBPF support, not a specific distro |
| `gore2regex` (loki.secretfilter, full profile) | Yes | Uses `wasilibs/go-re2`; built with `CGO_ENABLED=1` on EL8 |

`CGO_ENABLED=0` is **not** viable for the full profile (Oracle DB / journal / re2 paths).
The slim profile still uses `gore2regex` by default and keeps `CGO_ENABLED=1`.

## CI and local builds

GitHub Actions jobs use an `almalinux:8` container. Local parity:

```bash
./scripts/build-linux-amd64-container.sh slim
./scripts/build-linux-amd64-container.sh full
```

Verify the linked glibc ceiling after any build change:

```bash
VERIFY_GLIBC=1 ./scripts/build-linux-amd64-slim.sh
./scripts/verify-glibc-requirements.sh build/alloy
```

## Operational limits (not glibc-related)

- **Beyla / eBPF**: needs a recent kernel with BTF and appropriate capabilities; EL8
  kernels generally work when configured, but this is independent of the build host.
- **smartctl exporter**: requires `/usr/sbin/smartctl` (or configured path) on the host.
- **docker_state exporter**: requires access to the Docker socket or remote API.
