package smartctl

import (
	"time"

	"github.com/grafana/alloy/internal/component"
	"github.com/grafana/alloy/internal/component/prometheus/exporter"
	"github.com/grafana/alloy/internal/component/prometheus/exporter/common"
	"github.com/grafana/alloy/internal/featuregate"
	"github.com/grafana/alloy/internal/static/integrations"
	"github.com/grafana/alloy/internal/static/integrations/smartctl_exporter"
)

func init() {
	component.Register(component.Registration{
		Name:      "prometheus.exporter.smartctl",
		Stability: featuregate.StabilityGenerallyAvailable,
		Args:      Arguments{},
		Exports:   exporter.Exports{},

		Build: exporter.New(createIntegration, "smartctl"),
	})
}

func createIntegration(opts component.Options, args component.Arguments) (integrations.Integration, string, error) {
	common.WarningIfUsedInCluster(opts)
	a := args.(Arguments)
	defaultInstanceKey := common.HostNameInstanceKey()
	return integrations.NewIntegrationWithInstanceKey(opts.Logger, a.Convert(), defaultInstanceKey)
}

// DefaultArguments holds the default arguments for the prometheus.exporter.smartctl
// component.
var DefaultArguments = Arguments{
	SmartctlPath:     "/usr/sbin/smartctl",
	SmartctlInterval: 60 * time.Second,
	SmartctlRescan:   10 * time.Minute,
}

// Arguments configures the prometheus.exporter.smartctl component.
type Arguments struct {
	SmartctlPath     string        `alloy:"smartctl_path,attr,optional"`
	SmartctlInterval time.Duration `alloy:"smartctl_interval,attr,optional"`
	SmartctlRescan   time.Duration `alloy:"smartctl_rescan,attr,optional"`
	Devices          []string      `alloy:"devices,attr,optional"`
	DeviceExclude    string        `alloy:"device_exclude,attr,optional"`
	DeviceInclude    string        `alloy:"device_include,attr,optional"`
}

// SetToDefault implements syntax.Defaulter.
func (a *Arguments) SetToDefault() {
	*a = DefaultArguments
}

func (a *Arguments) Convert() *smartctl_exporter.Config {
	return &smartctl_exporter.Config{
		SmartctlPath:     a.SmartctlPath,
		SmartctlInterval: a.SmartctlInterval,
		SmartctlRescan:   a.SmartctlRescan,
		Devices:          a.Devices,
		DeviceExclude:    a.DeviceExclude,
		DeviceInclude:    a.DeviceInclude,
	}
}
