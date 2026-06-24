# sqlqueryreceiver build handoff (`v1.18.0-sql`)

Handoff for adding `otelcol.receiver.sqlquery` to the dontmint-alloy extended build for migrating central SQL exporter jobs to Alloy.

## What is already done in this branch

| Item | Location |
|------|----------|
| Alloy component wrapper | `internal/component/otelcol/receiver/sqlquery/` |
| Component registration | `internal/component/all/profile_*.go` |
| OTel collector distribution entry | `collector/builder-config.yaml` (full), `collector/builder-config.slim.yaml` (slim) |
| ClickHouse driver allow-list | `internal/opentelemetry/sqlquery/driver.go` (fork via `go.mod` replace) |
| Driver blank imports (CH + PG) | `internal/component/otelcol/receiver/sqlquery/sqlquery.go` |
| Unit test (unmarshal + convert) | `internal/component/otelcol/receiver/sqlquery/sqlquery_test.go` |
| Example Alloy job definitions | `infra-collab/monitoring/alloy/sqlquery_jobs.yml` (external repo) |

## Fork note: ClickHouse driver

Upstream `opentelemetry-collector-contrib` v0.147.0 `internal/sqlquery` does **not** list `clickhouse` in `IsValidDriver`, even though `clickhouse-go/v2` is already an indirect dependency of Alloy.

This build vendors a minimal fork:

```go
replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/sqlquery => ./internal/opentelemetry/sqlquery
```

Changes vs upstream: `DriverClickhouse = "clickhouse"` added to `driver.go`. Connections use existing `clickhouse://…` DSNs with `driver = "clickhouse"` and `datasource = "…"`.

Postgres continues to use `driver = "postgres"` with lib/pq (registered in upstream `sqlqueryreceiver/internal/database_sql.go` and duplicated in the Alloy wrapper for clarity).

## Build steps (do **not** run in investigation-only mode)

Full profile:

```bash
cd dontmint-alloy
export GO_TAGS='gore2regex embedalloyui'
./scripts/build-linux-amd64.sh
go test ./internal/component/otelcol/receiver/sqlquery/...
```

Slim profile (smaller binary, sqlquery included):

```bash
./scripts/build-linux-amd64-slim.sh
go test -tags alloy_slim ./internal/component/otelcol/receiver/sqlquery/...
```

See `docs/build-profiles.md` for profile details.

Verify the binary lists the component:

```bash
./build/alloy convert --source=otelcol --target=alloy <<'EOF'
receivers:
  sqlquery:
    driver: postgres
    datasource: "host=127.0.0.1 port=5432 user=u password=p sslmode=disable"
    queries:
      - sql: "select 1 as v"
        metrics:
          - metric_name: sql_example
            value_column: v
EOF
```

## Release tag

After CI build passes on the branch:

```bash
git tag v1.18.0-sql
git push origin v1.18.0-sql
```

GitHub Actions (`.github/workflows/alloy-build-and-release.yml`) publishes `alloy-linux-amd64.xz` on tag push `v*`.

## Runtime requirements

- Run Alloy with `--stability.level=public-preview` because `otelcol.receiver.sqlquery` is registered at **Public Preview** stability.
- Dual-run with existing SQL exporter until `sql_*` metric parity is confirmed in VictoriaMetrics; then retire the legacy scrape job.

## Known gaps (separate tracks)

| Gap | Notes |
|-----|-------|
| Cron scheduling | `AtomSync` (`55 * * * *`) uses `collection_interval = "60m"` as best-effort; sqlqueryreceiver has no cron. SQL is hour-anchored via `formatDateTime(now()-1h)`. |
| Zero-row queries | Verify `allow_zero_rows: true` queries (`coscan_*`, `agg_daily_view_is_empty`, merge/S3 metrics) still emit expected series during dual-run. |
| Multi-connection duplication | Legacy exporter runs kbf/kba/kbc separately; Alloy config mirrors that with one receiver per DSN. Alerts using `avg()`/`sum()` may differ slightly if series cardinality changes—validate during soak. |

## Config pipeline shape

```
otelcol.receiver.sqlquery "<job>_<conn>" { … }
  -> otelcol.exporter.prometheus "sql_exporter"
  -> prometheus.remote_write "victoriametrics"
```

Metric names use the `sql_<query_name>` prefix to match alert rule definitions in the monitoring repo.
