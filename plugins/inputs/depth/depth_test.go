package depth

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGathering(t *testing.T) {
	var depth = NewDepth()
	depth.Pairs = []string{"XRPUSDT"}
	depth.Count = 2

	var err error
	err = depth.Init()
	assert.NoError(t, err)

	if testing.Short() {
		t.Skip("Skipping network-dependent test in short mode.")
	}
	var acc testutil.Accumulator
	err = acc.GatherError(depth.Gather)
	assert.NoError(t, err)
	metric, ok := acc.Get("depth")
	require.True(t, ok)

	metricNames := []string{
		"asks_0_0",
		"asks_0_1",
		"asks_0_2",
		"asks_1_0",
		"asks_1_1",
		"asks_1_2",
		"bids_0_0",
		"bids_0_1",
		"bids_0_2",
		"bids_1_0",
		"bids_1_1",
		"bids_1_2",
	}
	for _, metricName := range metricNames {
		assert.Contains(t, metric.Fields, metricName)
		_, ok := metric.Fields[metricName].(float64)
		assert.True(t, true, ok)
	}
}
