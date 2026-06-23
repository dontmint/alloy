# sqlqueryreceiver build handoff (`v1.18.0-sql`)

Handoff for adding `otelcol.receiver.sqlquery` to the dontmint-alloy extended build so Scaleflex can migrate central `sql_exporter` jobs to Alloy.

## What is already done in this branch

| Item | Location |
|------|----------|
| Alloy component wrapper | `internal/component/otelcol/receiver/sqlquery/` |
| Component registration | `internal/component/all/all.go` |
| OTel collector distribution entry | `collector/builder-config.yaml` |
| ClickHouse driver allow-list | `internal/opentelemetry/sqlquery/driver.go` (fork via `go.mod` replace) |
| Driver blank imports (CH + PG) | `internal/component/otelcol/receiver/sqlquery/sqlquery.go` |
| Unit test (unmarshal + convert) | `internal/component/otelcol/receiver/sqlquery/sqlquery_test.go` |
| Ansible Alloy config translation | `scaleflex-ansible/roles/alloy/` + `infra-collab/monitoring/alloy/sqlquery_jobs.yml` |

## Fork note: ClickHouse driver

Upstream `opentelemetry-collector-contrib` v0.147.0 `internal/sqlquery` does **not** list `clickhouse` in `IsValidDriver`, even though `clickhouse-go/v2` is already an indirect dependency of Alloy.

This build vendors a minimal fork:

```go
replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/sqlquery => ./internal/opentelemetry/sqlquery
```

Changes vs upstream: `DriverClickhouse = "clickhouse"` added to `driver.go`. Connections use the existing `clickhouse://…` DSNs from `sql_exporter_connections` with `driver = "clickhouse"` and `datasource = "…"`.

Postgres continues to use `driver = "postgres"` with lib/pq (registered in upstream `sqlqueryreceiver/internal/database_sql.go` and duplicated in the Alloy wrapper for clarity).

## Build steps (do **not** run in investigation-only mode)

```bash
cd dontmint-alloy
export GO_TAGS='gore2regex embedalloyui'
./scripts/build-linux-amd64.sh
go test ./internal/component/otelcol/receiver/sqlquery/...
```

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

Update Scaleflex Ansible when cutting over:

```yaml
alloy_version: "v1.18.0-sql"
alloy_download_url_binary: "https://github.com/dontmint/alloy/releases/download/v1.18.0-sql/alloy-linux-amd64.xz"
```

## Runtime requirements

- Deploy on `metrics-vm-p001-fr-ov-lim1.metrics.scal3fl3x.com` (or dedicated host) with `--stability.level=public-preview` because `otelcol.receiver.sqlquery` is registered at **Public Preview** stability.
- Enable `alloy_integrations.sqlquery_receiver.enabled: true` in host vars (see `scaleflex-ansible/host_vars/metrics-vm-p001-fr-ov-lim1.metrics.scal3fl3x.com/alloy.yml`).
- Dual-run with existing `sql_exporter` until `sql_*` parity is confirmed in VictoriaMetrics; then retire the vmagent `sql_exporter` scrape job.

## Known gaps (separate tracks)

| Gap | Notes |
|-----|-------|
| Cron scheduling | `AtomSync` (`55 * * * *`) uses `collection_interval = "60m"` as best-effort; sqlqueryreceiver has no cron. SQL is hour-anchored via `formatDateTime(now()-1h)`. |
| Zero-row queries | Verify `allow_zero_rows: true` queries (`coscan_*`, `agg_daily_view_is_empty`, merge/S3 metrics) still emit expected series during dual-run. |
| Multi-connection duplication | sql_exporter runs kbf/kba/kbc separately; Alloy config mirrors that with one receiver per DSN. Alerts using `avg()`/`sum()` may differ slightly if series cardinality changes—validate during soak. |

## Config pipeline shape

```
otelcol.receiver.sqlquery "<job>_<conn>" { … }
  -> otelcol.exporter.prometheus "sql_exporter"
  -> prometheus.remote_write "victoriametrics"
```

Metric names use the `sql_<query_name>` prefix to match `infra-collab/monitoring/alerts/custom_rules_api.yml`.
