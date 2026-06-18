package docker_state_exporter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestDockerStateExporter_Config(t *testing.T) {
	var c Config

	err := yaml.Unmarshal([]byte("{}"), &c)
	require.NoError(t, err)
	require.Equal(t, DefaultConfig, c)
}

func TestDockerStateExporter_ConfigWithValues(t *testing.T) {
	var c Config

	err := yaml.Unmarshal([]byte(`
docker_host: tcp://192.168.1.1:2375
cache_period: 5
enable_labels: false
`), &c)
	require.NoError(t, err)
	require.Equal(t, "tcp://192.168.1.1:2375", c.DockerHost)
	require.Equal(t, 5, c.CachePeriod)
	require.False(t, c.EnableLabels)
}
