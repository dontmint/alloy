//go:build alloy_slim

// Slim build profile: host metrics, process/docker/smartctl/postgres exporters,
// Beyla eBPF, Prometheus scrape/relabel/remote_write, and otelcol.receiver.sqlquery.
package all

import (
	_ "github.com/grafana/alloy/internal/component/beyla/ebpf"                     // Import beyla.ebpf
	_ "github.com/grafana/alloy/internal/component/discovery/relabel"                // Import discovery.relabel
	_ "github.com/grafana/alloy/internal/component/local/file"                       // Import local.file
	_ "github.com/grafana/alloy/internal/component/local/file_match"                 // Import local.file_match
	_ "github.com/grafana/alloy/internal/component/otelcol/auth/basic"               // Import otelcol.auth.basic
	_ "github.com/grafana/alloy/internal/component/otelcol/exporter/debug"           // Import otelcol.exporter.debug
	_ "github.com/grafana/alloy/internal/component/otelcol/exporter/otlphttp"        // Import otelcol.exporter.otlphttp
	_ "github.com/grafana/alloy/internal/component/otelcol/exporter/prometheus"      // Import otelcol.exporter.prometheus
	_ "github.com/grafana/alloy/internal/component/otelcol/processor/attributes"     // Import otelcol.processor.attributes
	_ "github.com/grafana/alloy/internal/component/otelcol/processor/batch"          // Import otelcol.processor.batch
	_ "github.com/grafana/alloy/internal/component/otelcol/processor/memorylimiter"  // Import otelcol.processor.memory_limiter
	_ "github.com/grafana/alloy/internal/component/otelcol/receiver/sqlquery"        // Import otelcol.receiver.sqlquery
	_ "github.com/grafana/alloy/internal/component/otelcol/storage/file"             // Import otelcol.storage.file
	_ "github.com/grafana/alloy/internal/component/prometheus/exporter/docker_state" // Import prometheus.exporter.docker_state
	_ "github.com/grafana/alloy/internal/component/prometheus/exporter/postgres"     // Import prometheus.exporter.postgres
	_ "github.com/grafana/alloy/internal/component/prometheus/exporter/process"      // Import prometheus.exporter.process
	_ "github.com/grafana/alloy/internal/component/prometheus/exporter/self"         // Import prometheus.exporter.self
	_ "github.com/grafana/alloy/internal/component/prometheus/exporter/smartctl"     // Import prometheus.exporter.smartctl
	_ "github.com/grafana/alloy/internal/component/prometheus/exporter/static"       // Import prometheus.exporter.static
	_ "github.com/grafana/alloy/internal/component/prometheus/exporter/unix"         // Import prometheus.exporter.unix
	_ "github.com/grafana/alloy/internal/component/prometheus/relabel"               // Import prometheus.relabel
	_ "github.com/grafana/alloy/internal/component/prometheus/remotewrite"           // Import prometheus.remote_write
	_ "github.com/grafana/alloy/internal/component/prometheus/scrape"                // Import prometheus.scrape
)
