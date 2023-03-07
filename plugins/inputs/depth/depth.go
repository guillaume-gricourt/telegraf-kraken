//go:generate ../../../tools/readme_config_includer/generator
package depth

import (
	_ "embed"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/guillaume-gricourt/telegraf-kraken/pkg/krakenapi"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	jsonparser "github.com/influxdata/telegraf/plugins/parsers/json"
)

//go:embed depth.conf
var sampleConfig string

const urlSuffix = "/public/Depth"
const method = "GET"

type Depth struct {
	Pairs   []string      `toml:"pairs"`
	Count   int           `toml:"count"`
	Timeout time.Duration `toml:"timeout"`

	client *krakenapi.Client
}

func NewDepth() *Depth {
	return &Depth{
		Pairs:   []string{},
		Count:   100,
		Timeout: time.Second * 5,
	}
}

func (*Depth) SampleConfig() string {
	return sampleConfig
}

func (*Depth) Description() string {
	return "telegraf-kraken: depth"
}

func (d *Depth) Init() error {
	d.client = krakenapi.NewClient(method, urlSuffix, "", d.Timeout)
	if len(d.Pairs) < 1 {
		return errors.New("Provide at least one asset to download")
	}
	if d.Count < 1 || d.Count > 500 {
		return errors.New("count parameter must be comprise between 1 and 500")
	}
	return nil
}

func (d *Depth) Gather(accumulator telegraf.Accumulator) error {
	var err error
	// parameter
	parameters := map[string]string{"pair": strings.Join(d.Pairs, ","), "count": strconv.Itoa(d.Count)}
	// request
	resp, err := d.client.Request(nil, parameters)
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
		accumulator.AddFields("depth", flattener.Fields, map[string]string{"asset": pair})
	}
	return nil
}

func init() {
	inputs.Add("depth", func() telegraf.Input {
		return NewDepth()
	})
}
