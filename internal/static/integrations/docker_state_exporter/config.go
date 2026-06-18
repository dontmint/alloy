package docker_state_exporter

import (
	"log/slog"

	"github.com/grafana/alloy/internal/static/integrations"
	integrations_v2 "github.com/grafana/alloy/internal/static/integrations/v2"
	"github.com/grafana/alloy/internal/static/integrations/v2/metricsutils"
)

// DefaultConfig holds the default settings for the docker_state_exporter integration.
var DefaultConfig = Config{
	DockerHost:   "unix:///var/run/docker.sock",
	CachePeriod:  1,
	EnableLabels: true,
}

// Config controls the docker_state_exporter integration.
type Config struct {
	DockerHost   string `yaml:"docker_host,omitempty"`
	CachePeriod  int    `yaml:"cache_period,omitempty"`
	EnableLabels bool   `yaml:"enable_labels,omitempty"`
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (c *Config) UnmarshalYAML(unmarshal func(v any) error) error {
	*c = DefaultConfig

	type plain Config
	return unmarshal((*plain)(c))
}

// Name returns the name of the integration that this config represents.
func (c *Config) Name() string {
	return "docker_state_exporter"
}

func (c *Config) InstanceKey(defaultKey string) (string, error) {
	return defaultKey, nil
}

// NewIntegration converts this config into an instance of an integration.
func (c *Config) NewIntegration(l *slog.Logger) (integrations.Integration, error) {
	return New(l, c)
}

func init() {
	integrations.RegisterIntegration(&Config{})
	integrations_v2.RegisterLegacy(&Config{}, integrations_v2.TypeSingleton, metricsutils.NewNamedShim("docker_state"))
}
