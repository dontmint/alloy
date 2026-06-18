//go:build linux

package docker_state_exporter

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/grafana/alloy/internal/build"
	"github.com/grafana/alloy/internal/static/integrations/config"
)

// Integration is the docker_state_exporter integration. The integration scrapes
// metrics from Docker containers via the local Docker socket.
type Integration struct {
	c         *Config
	collector *dockerStateCollector
}

// New creates a new instance of the docker_state_exporter integration.
func New(_ *slog.Logger, c *Config) (*Integration, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithHost(c.DockerHost))
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	_, err = cli.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to docker daemon: %w", err)
	}

	return &Integration{
		c: c,
		collector: &dockerStateCollector{
			client:      cli,
			cachePeriod: time.Duration(c.CachePeriod) * time.Second,
			enableLabels: c.EnableLabels,
		},
	}, nil
}

// MetricsHandler satisfies Integration.RegisterRoutes.
func (i *Integration) MetricsHandler() (http.Handler, error) {
	r := prometheus.NewRegistry()
	if err := r.Register(i.collector); err != nil {
		return nil, fmt.Errorf("couldn't register docker_state_exporter collector: %w", err)
	}

	if err := r.Register(build.NewCollector("docker_state_exporter")); err != nil {
		return nil, fmt.Errorf("couldn't register docker_state_exporter: %w", err)
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
	<-ctx.Done()
	if i.collector.client != nil {
		i.collector.client.Close()
	}
	return ctx.Err()
}

// dockerStateCollector implements prometheus.Collector to expose Docker container state metrics.
type dockerStateCollector struct {
	client       *client.Client
	mu           sync.Mutex
	cachePeriod  time.Duration
	enableLabels bool
	lastSeen     time.Time
	cachedInfo   []container.InspectResponse
}

type descSource struct {
	name string
	help string
}

func (d *descSource) desc(labels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(d.name, d.help, nil, labels)
}

var (
	namespace = "container_state_"

	healthStatusDesc = &descSource{namespace + "health_status", "Container health status."}
	statusDesc       = &descSource{namespace + "status", "Container status."}
	oomKilledDesc    = &descSource{namespace + "oomkilled", "Container was killed by OOMKiller."}
	startedAtDesc    = &descSource{namespace + "startedat", "Time when the Container started."}
	finishedAtDesc   = &descSource{namespace + "finishedat", "Time when the Container finished."}
	restartCountDesc = &descSource{"container_restartcount", "Number of times the container has been restarted."}
	exitCodeDesc     = &descSource{"container_exitcode", "Exit code of the container."}
)

func (c *dockerStateCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- healthStatusDesc.desc(nil)
	ch <- statusDesc.desc(nil)
	ch <- oomKilledDesc.desc(nil)
	ch <- startedAtDesc.desc(nil)
	ch <- finishedAtDesc.desc(nil)
	ch <- restartCountDesc.desc(nil)
	ch <- exitCodeDesc.desc(nil)
}

func (c *dockerStateCollector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	if now.Sub(c.lastSeen) >= c.cachePeriod {
		_ = c.refreshCache()
		c.lastSeen = now
	}

	c.collectMetrics(ch)
}

func (c *dockerStateCollector) refreshCache() error {
	containers, err := c.client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return err
	}

	c.cachedInfo = nil
	for _, ctr := range containers {
		info, err := c.client.ContainerInspect(context.Background(), ctr.ID)
		if err != nil {
			continue
		}

		if info.Config == nil {
			info.Config = &container.Config{Labels: map[string]string{}}
		}
		if info.State == nil {
			info.State = &container.State{}
		}
		if info.State.Health == nil {
			info.State.Health = &container.Health{Status: container.NoHealthcheck}
		}

		c.cachedInfo = append(c.cachedInfo, info)
	}

	return nil
}

func (c *dockerStateCollector) collectMetrics(ch chan<- prometheus.Metric) {
	rep := regexp.MustCompile("[^a-zA-Z0-9_]")

	for _, info := range c.cachedInfo {
		labels := prometheus.Labels{}

		if c.enableLabels && info.Config != nil {
			for k, v := range info.Config.Labels {
				label := strings.ToLower("container_label_" + k)
				labels[rep.ReplaceAllLiteralString(label, "_")] = v
			}
		}

		labels["id"] = "/docker/" + info.ID
		labels["image"] = info.Config.Image
		labels["name"] = strings.TrimPrefix(info.Name, "/")

		boolToFloat := func(b bool) float64 {
			if b {
				return 1
			}
			return 0
		}

		copyLabels := func(src prometheus.Labels) prometheus.Labels {
			dst := make(prometheus.Labels, len(src))
			for k, v := range src {
				dst[k] = v
			}
			return dst
		}

		healthStatuses := []container.HealthStatus{container.NoHealthcheck, container.Starting, container.Healthy, container.Unhealthy}
		for _, status := range healthStatuses {
			tmpLabels := copyLabels(labels)
			tmpLabels["status"] = string(status)
			ch <- prometheus.MustNewConstMetric(healthStatusDesc.desc(tmpLabels), prometheus.GaugeValue, boolToFloat(info.State.Health.Status == status))
		}

		statuses := []string{"paused", "restarting", "running", "removing", "dead", "created", "exited"}
		for _, status := range statuses {
			tmpLabels := copyLabels(labels)
			tmpLabels["status"] = status
			ch <- prometheus.MustNewConstMetric(statusDesc.desc(tmpLabels), prometheus.GaugeValue, boolToFloat(string(info.State.Status) == status))
		}

		ch <- prometheus.MustNewConstMetric(oomKilledDesc.desc(labels), prometheus.GaugeValue, boolToFloat(info.State.OOMKilled))

		startedAt, _ := time.Parse(time.RFC3339Nano, info.State.StartedAt)
		ch <- prometheus.MustNewConstMetric(startedAtDesc.desc(labels), prometheus.GaugeValue, float64(startedAt.Unix()))

		finishedAt, _ := time.Parse(time.RFC3339Nano, info.State.FinishedAt)
		ch <- prometheus.MustNewConstMetric(finishedAtDesc.desc(labels), prometheus.GaugeValue, float64(finishedAt.Unix()))

		ch <- prometheus.MustNewConstMetric(restartCountDesc.desc(labels), prometheus.GaugeValue, float64(info.RestartCount))
		ch <- prometheus.MustNewConstMetric(exitCodeDesc.desc(labels), prometheus.GaugeValue, float64(info.State.ExitCode))
	}
}
