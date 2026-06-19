package build

import (
	"github.com/grafana/alloy/internal/component/discovery"
	"github.com/grafana/alloy/internal/component/prometheus/exporter/smartctl"
	"github.com/grafana/alloy/internal/static/integrations/smartctl_exporter"
)

func (b *ConfigBuilder) appendSmartctlExporter(config *smartctl_exporter.Config, instanceKey *string) discovery.Exports {
	args := toSmartctlExporter(config)
	return b.appendExporterBlock(args, config.Name(), instanceKey, "smartctl")
}

func toSmartctlExporter(config *smartctl_exporter.Config) *smartctl.Arguments {
	return &smartctl.Arguments{
		SmartctlPath:     config.SmartctlPath,
		SmartctlInterval: config.SmartctlInterval,
		SmartctlRescan:   config.SmartctlRescan,
		Devices:          config.Devices,
		DeviceExclude:    config.DeviceExclude,
		DeviceInclude:    config.DeviceInclude,
	}
}
