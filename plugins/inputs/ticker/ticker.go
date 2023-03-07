//go:generate ../../../tools/readme_config_includer/generator
package ticker

import (
	_ "embed"
	"errors"
	"strings"
	"time"

	"github.com/guillaume-gricourt/telegraf-kraken/pkg/krakenapi"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	jsonparser "github.com/influxdata/telegraf/plugins/parsers/json"
)

//go:embed sample.conf
var sampleConfig string

const urlSuffix = "/public/Ticker"
const method = "GET"

type Ticker struct {
	Pairs   []string      `toml:"pairs"`
	Timeout time.Duration `toml:"timeout"`

	client *krakenapi.Client
}

func NewTicker() *Ticker {
	return &Ticker{
		Pairs:   []string{},
		Timeout: time.Second * 5,
	}
}

func (*Ticker) SampleConfig() string {
	return sampleConfig
}

func (*Ticker) Description() string {
	return "telegraf-kraken: ticker"
}

func (t *Ticker) Init() error {
	t.client = krakenapi.NewClient(method, urlSuffix, "", t.Timeout)
	if len(t.Pairs) < 1 {
		return errors.New("Provide at least one asset to download")
	}
	return nil
}

func (t *Ticker) Gather(accumulator telegraf.Accumulator) error {
	var err error
	// parameter
	parameters := map[string]string{"pair": strings.Join(t.Pairs, ",")}
	// request
	resp, err := t.client.Request(nil, parameters)
	if err != nil {
		return err
	}
	// aggregate
	flattener := jsonparser.JSONFlattener{}
	for pair := range resp {
		err = flattener.FullFlattenJSON("", resp[pair], true, true)
		if err != nil {
			return err
		}
		accumulator.AddFields("ticker", flattener.Fields, map[string]string{"asset": pair})
	}
	return nil
}

func init() {
	inputs.Add("ticker", func() telegraf.Input {
		return NewTicker()
	})
}
