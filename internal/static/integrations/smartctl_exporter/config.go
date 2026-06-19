package smartctl_exporter

import (
	"log/slog"
	"time"

	"github.com/grafana/alloy/internal/static/integrations"
	integrations_v2 "github.com/grafana/alloy/internal/static/integrations/v2"
	"github.com/grafana/alloy/internal/static/integrations/v2/metricsutils"
)

// DefaultConfig holds the default settings for the smartctl_exporter integration.
var DefaultConfig = Config{
	SmartctlPath:     "/usr/sbin/smartctl",
	SmartctlInterval: 60 * time.Second,
	SmartctlRescan:   10 * time.Minute,
}

// Config controls the smartctl_exporter integration.
type Config struct {
	SmartctlPath     string        `yaml:"smartctl_path,omitempty"`
	SmartctlInterval time.Duration `yaml:"smartctl_interval,omitempty"`
	SmartctlRescan   time.Duration `yaml:"smartctl_rescan,omitempty"`
	Devices          []string      `yaml:"devices,omitempty"`
	DeviceExclude    string        `yaml:"device_exclude,omitempty"`
	DeviceInclude    string        `yaml:"device_include,omitempty"`
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (c *Config) UnmarshalYAML(unmarshal func(v any) error) error {
	*c = DefaultConfig

	type plain Config
	return unmarshal((*plain)(c))
}

// Name returns the name of the integration that this config represents.
func (c *Config) Name() string {
	return "smartctl_exporter"
}

func (c *Config) InstanceKey(defaultKey string) (string, error) {
	return defaultKey, nil
}

// NewIntegration converts this config into an instance of an integration.
func (c *Config) NewIntegration(l *slog.Logger) (integrations.Integration, error) {
	return New(l, c)
}

func (c *Config) runtime() *Runtime {
	r := &Runtime{
		Path:          c.SmartctlPath,
		Interval:      c.SmartctlInterval,
		DeviceExclude: c.DeviceExclude,
		DeviceInclude: c.DeviceInclude,
	}
	r.initCache()
	return r
}

func init() {
	integrations.RegisterIntegration(&Config{})
	integrations_v2.RegisterLegacy(&Config{}, integrations_v2.TypeSingleton, metricsutils.NewNamedShim("smartctl"))
}
