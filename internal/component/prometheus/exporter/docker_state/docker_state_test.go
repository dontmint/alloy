package docker_state

import (
	"testing"

	"github.com/grafana/alloy/syntax"
	"github.com/stretchr/testify/require"
)

func TestAlloyConfigUnmarshal(t *testing.T) {
	var exampleAlloyConfig = `
	docker_host = "tcp://192.168.1.1:2375"
	cache_period = 5
	enable_labels = false
`

	var args Arguments
	err := syntax.Unmarshal([]byte(exampleAlloyConfig), &args)
	require.NoError(t, err)

	require.Equal(t, "tcp://192.168.1.1:2375", args.DockerHost)
	require.Equal(t, 5, args.CachePeriod)
	require.False(t, args.EnableLabels)
}

func TestAlloyConfigConvert(t *testing.T) {
	var exampleAlloyConfig = `
	docker_host = "tcp://192.168.1.1:2375"
	cache_period = 5
	enable_labels = false
`

	var args Arguments
	err := syntax.Unmarshal([]byte(exampleAlloyConfig), &args)
	require.NoError(t, err)

	c := args.Convert()
	require.Equal(t, "tcp://192.168.1.1:2375", c.DockerHost)
	require.Equal(t, 5, c.CachePeriod)
	require.False(t, c.EnableLabels)
}
