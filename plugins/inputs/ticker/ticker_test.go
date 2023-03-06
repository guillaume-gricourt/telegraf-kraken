package ticker

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGathering(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping network-dependent test in short mode.")
	}
	var ticker = NewTicker()
	ticker.Pairs = []string{"XRPUSDT"}

	var err error
	err = ticker.Init()
	assert.NoError(t, err)

	var acc testutil.Accumulator
	err = acc.GatherError(ticker.Gather)
	assert.NoError(t, err)
	metric, ok := acc.Get("ticker")
	require.True(t, ok)

	metricNames := []string{
		"a_0",
		"a_1",
		"a_2",
		"b_0",
		"b_1",
		"b_2",
		"c_0",
		"c_1",
		"h_0",
		"h_1",
		"l_0",
		"l_1",
		"o",
		"p_0",
		"p_1",
		"t_0",
		"t_1",
		"v_0",
		"v_1",
	}
	for _, metricName := range metricNames {
		assert.Contains(t, metric.Fields, metricName)
		_, ok := metric.Fields[metricName].(float64)
		assert.True(t, true, ok)
	}
}
