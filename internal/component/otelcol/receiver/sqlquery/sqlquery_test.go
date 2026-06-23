package sqlquery_test

import (
	"testing"
	"time"

	"github.com/grafana/alloy/internal/component/otelcol/receiver/sqlquery"
	"github.com/grafana/alloy/syntax"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlqueryreceiver"
	"github.com/stretchr/testify/require"
)

func TestArguments_UnmarshalAlloy(t *testing.T) {
	cfg := `
driver = "clickhouse"
datasource = "clickhouse://user:pass@localhost:9000/default"
collection_interval = "60s"

query {
	sql = "SELECT 1 AS value"
	metric {
		metric_name = "sql_example"
		value_column = "value"
		data_type = "gauge"
		value_type = "double"
	}
}

output {}
`

	var args sqlquery.Arguments
	require.NoError(t, syntax.Unmarshal([]byte(cfg), &args))

	actualPtr, err := args.Convert()
	require.NoError(t, err)

	actual := actualPtr.(*sqlqueryreceiver.Config)
	require.Equal(t, "clickhouse", actual.Driver)
	require.Equal(t, "clickhouse://user:pass@localhost:9000/default", actual.DataSource)
	require.Equal(t, 60*time.Second, actual.CollectionInterval)
	require.Len(t, actual.Queries, 1)
	require.Equal(t, "SELECT 1 AS value", actual.Queries[0].SQL)
	require.Len(t, actual.Queries[0].Metrics, 1)
	require.Equal(t, "sql_example", actual.Queries[0].Metrics[0].MetricName)
	require.Equal(t, "value", actual.Queries[0].Metrics[0].ValueColumn)
	require.Equal(t, "gauge", string(actual.Queries[0].Metrics[0].DataType))
	require.Equal(t, "double", string(actual.Queries[0].Metrics[0].ValueType))
}
