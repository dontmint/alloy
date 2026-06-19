package smartctl

import (
	"testing"
	"time"

	"github.com/grafana/alloy/syntax"
	"github.com/stretchr/testify/require"
)

func TestAlloyConfigUnmarshal(t *testing.T) {
	var exampleAlloyConfig = `
	smartctl_path = "/usr/bin/smartctl"
	smartctl_interval = "30s"
	smartctl_rescan = "5m"
	devices = ["/dev/sda", "/dev/sdb"]
	device_exclude = "sd[a-b]"
	device_include = "nvme.*"
`

	var args Arguments
	err := syntax.Unmarshal([]byte(exampleAlloyConfig), &args)
	require.NoError(t, err)

	require.Equal(t, "/usr/bin/smartctl", args.SmartctlPath)
	require.Equal(t, 30*time.Second, args.SmartctlInterval)
	require.Equal(t, 5*time.Minute, args.SmartctlRescan)
	require.Equal(t, []string{"/dev/sda", "/dev/sdb"}, args.Devices)
	require.Equal(t, "sd[a-b]", args.DeviceExclude)
	require.Equal(t, "nvme.*", args.DeviceInclude)
}

func TestAlloyConfigConvert(t *testing.T) {
	var exampleAlloyConfig = `
	smartctl_path = "/usr/bin/smartctl"
	smartctl_interval = "30s"
	smartctl_rescan = "5m"
`

	var args Arguments
	err := syntax.Unmarshal([]byte(exampleAlloyConfig), &args)
	require.NoError(t, err)

	c := args.Convert()
	require.Equal(t, "/usr/bin/smartctl", c.SmartctlPath)
	require.Equal(t, 30*time.Second, c.SmartctlInterval)
	require.Equal(t, 5*time.Minute, c.SmartctlRescan)
}
