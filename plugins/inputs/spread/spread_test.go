package spread

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGathering(t *testing.T) {
	var spread = NewSpread()
	spread.Pairs = []string{"XRPUSDT"}

	var err error
	err = spread.Init()
	assert.NoError(t, err)

	if testing.Short() {
		t.Skip("Skipping network-dependent test in short mode.")
	}
	var acc testutil.Accumulator
	err = acc.GatherError(spread.Gather)
	assert.NoError(t, err)
	metric, ok := acc.Get("spread")
	require.True(t, ok)

	assert.Equal(t, metric.Measurement, "spread")
	metricNames := []string{
		"a",
		"b",
		"t",
	}
	for _, metricName := range metricNames {
		assert.Contains(t, metric.Fields, metricName)
		_, ok := metric.Fields[metricName].(float64)
		assert.True(t, true, ok)
	}
}
