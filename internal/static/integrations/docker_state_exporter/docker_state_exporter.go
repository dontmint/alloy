//go:build !linux

package docker_state_exporter

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/grafana/alloy/internal/static/integrations/config"
)

// Integration is the docker_state_exporter integration. On non-Linux platforms,
// this integration does nothing and will print a warning if enabled.
type Integration struct {
	c *Config
}

// New creates a docker_state_exporter integration for non-Linux platforms, which is always a
// no-op.
func New(logger *slog.Logger, c *Config) (*Integration, error) {
	logger.Warn("the docker_state_exporter only works on Linux; enabling it otherwise will do nothing")
	return &Integration{c: c}, nil
}

// MetricsHandler satisfies Integration.RegisterRoutes.
func (i *Integration) MetricsHandler() (http.Handler, error) {
	return http.NotFoundHandler(), nil
}

// ScrapeConfigs satisfies Integration.ScrapeConfigs.
func (i *Integration) ScrapeConfigs() []config.ScrapeConfig {
	return []config.ScrapeConfig{}
}

// Run satisfies Integration.Run.
func (i *Integration) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}
