//go:generate ../../../tools/readme_config_includer/generator
package ticker

import (
	_ "embed"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/guillaume-gricourt/telegraf-kraken/pkg/krakenapi"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	jsonparser "github.com/influxdata/telegraf/plugins/parsers/json"
)

//go:embed ticker.conf
var sampleConfig string

const urlSuffix = "/public/Ticker"
const method = "GET"

type Data struct {
	A []string `json:"a"`
	B []string `json:"b"`
	C []string `json:"c"`
	V []string `json:"v"`
	P []string `json:"p"`
	T []int    `json:"t"`
	L []string `json:"l"`
	H []string `json:"h"`
	O string   `json:"o"`
}

type Response struct {
	Result map[string]Data `json:"result"`
	Error  []string        `json:"error"`
}

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
	resp := &Response{}
	// parameter
	parameters := map[string]string{"pair": strings.Join(t.Pairs, ",")}
	// request
	err = t.client.Request(nil, parameters, resp)
	if err != nil {
		return err
	}
	if len(resp.Error) > 0 {
		return errors.New(strings.Join(resp.Error, ","))
	}
	// aggregate
	for pair := range resp.Result {
		var record map[string]interface{}
		jrec, err := json.Marshal(resp.Result[pair])
		if err != nil {
			return err
		}
		err = json.Unmarshal(jrec, &record)
		if err != nil {
			return err
		}
		flattener := jsonparser.JSONFlattener{}
		err = flattener.FullFlattenJSON("", record, true, true)
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
