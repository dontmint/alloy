package build

import (
	"github.com/grafana/alloy/internal/component/discovery"
	"github.com/grafana/alloy/internal/component/prometheus/exporter/docker_state"
	"github.com/grafana/alloy/internal/static/integrations/docker_state_exporter"
)

func (b *ConfigBuilder) appendDockerStateExporter(config *docker_state_exporter.Config, instanceKey *string) discovery.Exports {
	args := toDockerStateExporter(config)
	return b.appendExporterBlock(args, config.Name(), instanceKey, "docker_state")
}

func toDockerStateExporter(config *docker_state_exporter.Config) *docker_state.Arguments {
	return &docker_state.Arguments{
		DockerHost:   config.DockerHost,
		CachePeriod:  config.CachePeriod,
		EnableLabels: config.EnableLabels,
	}
}
