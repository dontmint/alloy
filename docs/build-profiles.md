# Build profiles

Alloy supports compile-time build profiles that control which Flow components,
static integrations, embedded UI assets, and OpenTelemetry Collector modules are
linked into the binary.

Select a profile with `ALLOY_BUILD_PROFILE` when invoking `make alloy` or the
helper scripts under `scripts/`.

## Profiles

| Profile | Tag | UI | OTel builder config |
| --- | --- | --- | --- |
| `full` (default) | _(none)_ | embedded (`embedalloyui`) | `collector/builder-config.yaml` |
| `slim` | `alloy_slim` | none (`SKIP_UI_BUILD=1`) | `collector/builder-config.slim.yaml` |

## Full profile

The default upstream-style distribution. All Flow components, vendor Prometheus
exporters, Kubernetes operators, Loki/Pyroscope components, and the embedded
web UI are included.

```bash
./scripts/build-linux-amd64.sh
# or
make alloy
```

## Slim profile

A metrics-focused distribution for host monitoring, Beyla eBPF, and SQL query metrics:

**Flow components**

- `beyla.ebpf`
- `prometheus.exporter.unix`, `process`, `docker_state`, `smartctl`, `postgres`, `self`, `static`
- `prometheus.scrape`, `prometheus.relabel`, `prometheus.remote_write`
- `discovery.relabel`, `local.file`, `local.file_match`
- `otelcol.receiver.sqlquery`, `otelcol.exporter.prometheus`
- `otelcol.processor.batch`, `otelcol.processor.attributes`, `otelcol.processor.memorylimiter`
- `otelcol.exporter.otlphttp`, `otelcol.auth.basic`, `otelcol.storage.file`, `otelcol.exporter.debug`

**Excluded**

- Embedded web UI (`embedalloyui`)
- Kubernetes Prometheus operators (`prometheus.operator.*`)
- Vendor Prometheus exporters (Apache, SNMP, cloud vendor exporters, etc.)
- Loki, Pyroscope, Faro, database observability
- Most OTel receivers/exporters beyond sqlquery, Beyla traces, and prometheus remote write

```bash
./scripts/build-linux-amd64-slim.sh
# or
ALLOY_BUILD_PROFILE=slim make alloy
```

## Build-time toggling

Profiles are implemented with Go build tags:

- `internal/component/all/profile_full.go` — `//go:build !alloy_slim`
- `internal/component/all/profile_slim.go` — `//go:build alloy_slim`
- `internal/static/integrations/install/profile_*.go` — same split

The Makefile sets `GO_TAGS` automatically:

- **full:** prepends `gore2regex`, keeps `embedalloyui` when set by scripts
- **slim:** adds `alloy_slim`, strips `embedalloyui`, sets `SKIP_UI_BUILD=1`

OTel collector codegen reads `BUILDER_CONFIG` (set from `ALLOY_BUILD_PROFILE`).
Regenerate the collector distro whenever you switch profiles before building:

```bash
ALLOY_BUILD_PROFILE=slim make generate-otel-collector-distro
ALLOY_BUILD_PROFILE=slim make alloy
```

## Adding components to slim

1. Add the blank import to `internal/component/all/profile_slim.go`.
2. If the component needs a static integration, add it to
   `internal/static/integrations/install/profile_slim.go`.
3. If the component wraps an OTel module, add the module to
   `collector/builder-config.slim.yaml` and regenerate the collector distro.

Keep slim intentionally small; prefer the full profile when you need logging,
tracing, Kubernetes operators, or vendor exporters.

Release binaries are built on AlmaLinux 8 (glibc 2.28). See
`docs/linux-binary-compatibility.md` for supported distros.
