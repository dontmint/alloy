//go:build alloy_slim

package install

import (
	_ "github.com/grafana/alloy/internal/static/integrations/docker_state_exporter" // register docker_state_exporter
	_ "github.com/grafana/alloy/internal/static/integrations/node_exporter"          // register node_exporter
	_ "github.com/grafana/alloy/internal/static/integrations/postgres_exporter"      // register postgres_exporter
	_ "github.com/grafana/alloy/internal/static/integrations/process_exporter"       // register process_exporter
	_ "github.com/grafana/alloy/internal/static/integrations/smartctl_exporter"     // register smartctl_exporter
)
