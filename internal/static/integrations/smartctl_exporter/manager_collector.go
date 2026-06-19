// Copyright 2022 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package smartctl_exporter

import (
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// ManagerCollector implements the Collector interface for smartctl_exporter.
type ManagerCollector struct {
	logger  log.Logger
	runtime *Runtime
	devices []string

	mutex sync.Mutex
}

// Describe sends the super-set of all possible descriptors of metrics
func (c *ManagerCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *ManagerCollector) Collect(ch chan<- prometheus.Metric) {
	info := NewSMARTctlInfo(ch)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, device := range c.devices {
		json := c.runtime.readData(c.logger, device)
		if json.Exists() {
			info.SetJSON(json)
			smart := NewSMARTctl(c.logger, json, ch)
			smart.Collect()
		}
	}
	ch <- prometheus.MustNewConstMetric(
		metricDeviceCount,
		prometheus.GaugeValue,
		float64(len(c.devices)),
	)
	info.Collect()
}

func (c *ManagerCollector) rescanForDevices(rescanInterval time.Duration) {
	for {
		time.Sleep(rescanInterval)
		level.Info(c.logger).Log("msg", "Rescanning for devices")
		devices := c.runtime.scanDevices(c.logger)
		c.mutex.Lock()
		c.devices = devices
		c.mutex.Unlock()
	}
}
