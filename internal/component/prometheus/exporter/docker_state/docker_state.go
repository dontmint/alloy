package docker_state

import (
	"github.com/grafana/alloy/internal/component"
	"github.com/grafana/alloy/internal/component/prometheus/exporter"
	"github.com/grafana/alloy/internal/component/prometheus/exporter/common"
	"github.com/grafana/alloy/internal/featuregate"
	"github.com/grafana/alloy/internal/static/integrations"
	"github.com/grafana/alloy/internal/static/integrations/docker_state_exporter"
)

func init() {
	component.Register(component.Registration{
		Name:      "prometheus.exporter.docker_state",
		Stability: featuregate.StabilityGenerallyAvailable,
		Args:      Arguments{},
		Exports:   exporter.Exports{},

		Build: exporter.New(createIntegration, "docker_state"),
	})
}

func createIntegration(opts component.Options, args component.Arguments) (integrations.Integration, string, error) {
	common.WarningIfUsedInCluster(opts)
	a := args.(Arguments)
	defaultInstanceKey := common.HostNameInstanceKey()
	return integrations.NewIntegrationWithInstanceKey(opts.Logger, a.Convert(), defaultInstanceKey)
}

// DefaultArguments holds the default arguments for the prometheus.exporter.docker_state
// component.
var DefaultArguments = Arguments{
	DockerHost:   "unix:///var/run/docker.sock",
	CachePeriod:  1,
	EnableLabels: true,
}

// Arguments configures the prometheus.exporter.docker_state component.
type Arguments struct {
	DockerHost   string `alloy:"docker_host,attr,optional"`
	CachePeriod  int    `alloy:"cache_period,attr,optional"`
	EnableLabels bool   `alloy:"enable_labels,attr,optional"`
}

// SetToDefault implements syntax.Defaulter.
func (a *Arguments) SetToDefault() {
	*a = DefaultArguments
}

func (a *Arguments) Convert() *docker_state_exporter.Config {
	return &docker_state_exporter.Config{
		DockerHost:   a.DockerHost,
		CachePeriod:  a.CachePeriod,
		EnableLabels: a.EnableLabels,
	}
}
