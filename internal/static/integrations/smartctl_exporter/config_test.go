package smartctl_exporter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestSmartctlExporter_Config(t *testing.T) {
	var c Config

	err := yaml.Unmarshal([]byte("{}"), &c)
	require.NoError(t, err)
	require.Equal(t, DefaultConfig, c)
}

func TestSmartctlExporter_ConfigWithValues(t *testing.T) {
	var c Config

	err := yaml.Unmarshal([]byte(`
smartctl_path: /usr/bin/smartctl
smartctl_interval: 30s
smartctl_rescan: 5m
devices:
  - /dev/sda
device_exclude: sd[a-b]
device_include: nvme.*
`), &c)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/smartctl", c.SmartctlPath)
	require.Equal(t, 30*time.Second, c.SmartctlInterval)
	require.Equal(t, 5*time.Minute, c.SmartctlRescan)
	require.Equal(t, []string{"/dev/sda"}, c.Devices)
	require.Equal(t, "sd[a-b]", c.DeviceExclude)
	require.Equal(t, "nvme.*", c.DeviceInclude)
}
