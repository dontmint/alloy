//go:build linux

package smartctl_exporter

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/grafana/alloy/internal/build"
	"github.com/grafana/alloy/internal/slogadapter"
	"github.com/grafana/alloy/internal/static/integrations/config"
)

// Integration is the smartctl_exporter integration.
type Integration struct {
	c         *Config
	collector *ManagerCollector
}

// New creates a new instance of the smartctl_exporter integration.
func New(logger *slog.Logger, c *Config) (*Integration, error) {
	gkLogger := log.With(slogadapter.GoKit(logger.Handler()), "integration", c.Name())

	runtime := c.runtime()
	devices := c.Devices
	if len(devices) == 0 {
		level.Info(gkLogger).Log("msg", "No devices specified, trying to load them automatically")
		devices = runtime.scanDevices(gkLogger)
		level.Info(gkLogger).Log("msg", "Number of devices found", "count", len(devices))
	}

	collector := &ManagerCollector{
		logger:  gkLogger,
		runtime: runtime,
		devices: devices,
	}

	return &Integration{
		c:         c,
		collector: collector,
	}, nil
}

// MetricsHandler satisfies Integration.RegisterRoutes.
func (i *Integration) MetricsHandler() (http.Handler, error) {
	r := prometheus.NewRegistry()
	if err := r.Register(i.collector); err != nil {
		return nil, fmt.Errorf("couldn't register smartctl_exporter collector: %w", err)
	}

	if err := r.Register(build.NewCollector("smartctl_exporter")); err != nil {
		return nil, fmt.Errorf("couldn't register smartctl_exporter: %w", err)
	}

	return promhttp.HandlerFor(
		r,
		promhttp.HandlerOpts{
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: 0,
		},
	), nil
}

// ScrapeConfigs satisfies Integration.ScrapeConfigs.
func (i *Integration) ScrapeConfigs() []config.ScrapeConfig {
	return []config.ScrapeConfig{{
		JobName:     i.c.Name(),
		MetricsPath: "/metrics",
	}}
}

// Run satisfies Integration.Run.
func (i *Integration) Run(ctx context.Context) error {
	if i.c.SmartctlRescan >= time.Second && len(i.c.Devices) == 0 {
		level.Info(i.collector.logger).Log("msg", "Start background scan process")
		level.Info(i.collector.logger).Log("msg", "Rescanning for devices every", "rescanInterval", i.c.SmartctlRescan)
		go i.collector.rescanForDevices(i.c.SmartctlRescan)
	}

	<-ctx.Done()
	return ctx.Err()
}
