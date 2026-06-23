package sqlquery

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// MetricArguments configures one OTel metric emitted from a SQL row.
type MetricArguments struct {
	MetricName       string            `alloy:"metric_name,attr" mapstructure:"metric_name"`
	ValueColumn      string            `alloy:"value_column,attr" mapstructure:"value_column"`
	AttributeColumns []string          `alloy:"attribute_columns,attr,optional" mapstructure:"attribute_columns"`
	DataType         string            `alloy:"data_type,attr,optional" mapstructure:"data_type"`
	ValueType        string            `alloy:"value_type,attr,optional" mapstructure:"value_type"`
	Monotonic        bool              `alloy:"monotonic,attr,optional" mapstructure:"monotonic"`
	Aggregation      string            `alloy:"aggregation,attr,optional" mapstructure:"aggregation"`
	Unit             string            `alloy:"unit,attr,optional" mapstructure:"unit"`
	Description      string            `alloy:"description,attr,optional" mapstructure:"description"`
	StaticAttributes map[string]string `alloy:"static_attributes,attr,optional" mapstructure:"static_attributes"`
	StartTsColumn    string            `alloy:"start_ts_column,attr,optional" mapstructure:"start_ts_column"`
	TsColumn         string            `alloy:"ts_column,attr,optional" mapstructure:"ts_column"`
}

// TelemetryArguments configures receiver self-telemetry.
type TelemetryArguments struct {
	Logs TelemetryLogsArguments `alloy:"logs,block,optional" mapstructure:"logs"`
}

// TelemetryLogsArguments configures query logging.
type TelemetryLogsArguments struct {
	Query bool `alloy:"query,attr,optional" mapstructure:"query"`
}

// QueryArguments configures one SQL statement and its metric mappings.
type QueryArguments struct {
	SQL                string            `alloy:"sql,attr" mapstructure:"sql"`
	Metrics            []MetricArguments `alloy:"metric,block,optional" mapstructure:"metrics"`
	TrackingColumn     string            `alloy:"tracking_column,attr,optional" mapstructure:"tracking_column"`
	TrackingStartValue string            `alloy:"tracking_start_value,attr,optional" mapstructure:"tracking_start_value"`
}

type argumentsDTO struct {
	Driver           string             `mapstructure:"driver"`
	DataSource       string             `mapstructure:"datasource"`
	Host             string             `mapstructure:"host"`
	Port             int                `mapstructure:"port"`
	Database         string             `mapstructure:"database"`
	Username         string             `mapstructure:"username"`
	Password         string             `mapstructure:"password"`
	AdditionalParams map[string]any     `mapstructure:"additional_params"`
	MaxOpenConn      int                `mapstructure:"max_open_conn"`
	Queries          []QueryArguments   `mapstructure:"queries"`
	Telemetry        TelemetryArguments `mapstructure:"telemetry"`
}

func (args Arguments) toDTO() (argumentsDTO, error) {
	if args.Driver == "" {
		return argumentsDTO{}, fmt.Errorf("driver must not be empty")
	}
	if len(args.Queries) == 0 {
		return argumentsDTO{}, fmt.Errorf("at least one query block is required")
	}

	for i, query := range args.Queries {
		if query.SQL == "" {
			return argumentsDTO{}, fmt.Errorf("query block %d: sql must not be empty", i)
		}
		if len(query.Metrics) == 0 {
			return argumentsDTO{}, fmt.Errorf("query block %d: at least one metric block is required", i)
		}
	}

	return argumentsDTO{
		Driver:           args.Driver,
		DataSource:       args.DataSource,
		Host:             args.Host,
		Port:             args.Port,
		Database:         args.Database,
		Username:         args.Username,
		Password:         args.Password,
		AdditionalParams: args.AdditionalParams,
		MaxOpenConn:      args.MaxOpenConn,
		Queries:          args.Queries,
		Telemetry:        args.Telemetry,
	}, nil
}

func decodeArguments(args Arguments, cfg any) error {
	dto, err := args.toDTO()
	if err != nil {
		return err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "mapstructure",
		Result:           cfg,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return fmt.Errorf("creating decoder: %w", err)
	}
	if err := decoder.Decode(dto); err != nil {
		return fmt.Errorf("decoding sqlquery config: %w", err)
	}
	return nil
}
