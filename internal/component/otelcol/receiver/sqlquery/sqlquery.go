// Package sqlquery provides an otelcol.receiver.sqlquery component.
package sqlquery

import (
	"fmt"

	"github.com/grafana/alloy/internal/component"
	"github.com/grafana/alloy/internal/component/otelcol"
	otelcolCfg "github.com/grafana/alloy/internal/component/otelcol/config"
	"github.com/grafana/alloy/internal/component/otelcol/extension"
	"github.com/grafana/alloy/internal/component/otelcol/receiver"
	"github.com/grafana/alloy/internal/featuregate"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlqueryreceiver"
	otelcomponent "go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pipeline"

	// Register SQL drivers used by Scaleflex sql_exporter parity (ClickHouse + Postgres).
	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/lib/pq"
)

func init() {
	component.Register(component.Registration{
		Name:      "otelcol.receiver.sqlquery",
		Stability: featuregate.StabilityPublicPreview,
		Args:      Arguments{},

		Build: func(opts component.Options, args component.Arguments) (component.Component, error) {
			fact := sqlqueryreceiver.NewFactory()
			return receiver.New(opts, fact, args.(Arguments))
		},
	})
}

// Arguments configures the otelcol.receiver.sqlquery component.
type Arguments struct {
	ScraperController otelcol.ScraperControllerArguments `alloy:",squash"`

	Driver           string         `alloy:"driver,attr"`
	DataSource       string         `alloy:"datasource,attr,optional"`
	Host             string         `alloy:"host,attr,optional"`
	Port             int            `alloy:"port,attr,optional"`
	Database         string         `alloy:"database,attr,optional"`
	Username         string         `alloy:"username,attr,optional"`
	Password         string         `alloy:"password,attr,optional"`
	AdditionalParams map[string]any `alloy:"additional_params,attr,optional"`
	MaxOpenConn      int            `alloy:"max_open_conn,attr,optional"`

	Queries   []QueryArguments   `alloy:"query,block,optional"`
	Telemetry TelemetryArguments `alloy:"telemetry,block,optional"`

	Storage *extension.ExtensionHandler `alloy:"storage,attr,optional"`

	DebugMetrics otelcolCfg.DebugMetricsArguments `alloy:"debug_metrics,block,optional"`

	Output *otelcol.ConsumerArguments `alloy:"output,block"`
}

var _ receiver.Arguments = Arguments{}

// SetToDefault implements syntax.Defaulter.
func (args *Arguments) SetToDefault() {
	def := sqlqueryreceiver.NewFactory().CreateDefaultConfig().(*sqlqueryreceiver.Config)
	*args = Arguments{
		ScraperController: otelcol.ScraperControllerArguments{
			CollectionInterval: def.CollectionInterval,
			InitialDelay:       def.InitialDelay,
			Timeout:            def.Timeout,
		},
		Output: &otelcol.ConsumerArguments{},
	}
	args.DebugMetrics.SetToDefault()
}

// Convert implements receiver.Arguments.
func (args Arguments) Convert() (otelcomponent.Config, error) {
	cfg := sqlqueryreceiver.NewFactory().CreateDefaultConfig().(*sqlqueryreceiver.Config)
	if controller := args.ScraperController.Convert(); controller != nil {
		cfg.ControllerConfig = *controller
	}
	if err := decodeArguments(args, cfg); err != nil {
		return nil, err
	}
	if args.Storage != nil {
		if args.Storage.Extension == nil {
			return nil, fmt.Errorf("storage extension %q is not running", args.Storage.ID)
		}
		cfg.StorageID = &args.Storage.ID
	}
	return cfg, nil
}

// Extensions implements receiver.Arguments.
func (args Arguments) Extensions() map[otelcomponent.ID]otelcomponent.Component {
	if args.Storage == nil || args.Storage.Extension == nil {
		return nil
	}
	return map[otelcomponent.ID]otelcomponent.Component{
		args.Storage.ID: args.Storage.Extension,
	}
}

// Exporters implements receiver.Arguments.
func (args Arguments) Exporters() map[pipeline.Signal]map[otelcomponent.ID]otelcomponent.Component {
	return nil
}

// NextConsumers implements receiver.Arguments.
func (args Arguments) NextConsumers() *otelcol.ConsumerArguments {
	return args.Output
}

// DebugMetricsConfig implements receiver.Arguments.
func (args Arguments) DebugMetricsConfig() otelcolCfg.DebugMetricsArguments {
	return args.DebugMetrics
}
