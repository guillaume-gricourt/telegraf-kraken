//go:generate ../../../tools/readme_config_includer/generator
package spread

import (
	_ "embed"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/guillaume-gricourt/telegraf-kraken/pkg/krakenapi"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//go:embed spread.conf
var sampleConfig string

const urlSuffix = "/public/Spread"
const method = "GET"

type Spread struct {
	Pairs   []string      `toml:"pairs"`
	Timeout time.Duration `toml:"timeout"`

	client *krakenapi.Client
}

func NewSpread() *Spread {
	return &Spread{
		Pairs:   []string{},
		Timeout: time.Second * 5,
	}
}

func (*Spread) SampleConfig() string {
	return sampleConfig
}

func (*Spread) Description() string {
	return "telegraf-kraken: spread"
}

func (s *Spread) Init() error {
	s.client = krakenapi.NewClient(method, urlSuffix, "", s.Timeout)
	if len(s.Pairs) < 1 {
		return errors.New("Provide at least one asset to download")
	}
	return nil
}

func (s *Spread) Gather(accumulator telegraf.Accumulator) error {
	var err error
	// parameter
	parameters := map[string]string{"pair": strings.Join(s.Pairs, ",")}
	// request
	resp, err := s.client.Request(nil, parameters)
	if err != nil {
		return err
	}
	// aggregate
	last := int64(resp["last"].(float64))
	delete(resp, "last")
	for pair := range resp {
		data := make(map[string]interface{})
		for _, value := range resp[pair].([]interface{}) {
			data["t"] = int64(value.([]interface{})[0].(float64))
			if data["t"] != last {
				continue
			}
			data["a"], _ = strconv.ParseFloat(value.([]interface{})[1].(string), 64)
			data["b"], _ = strconv.ParseFloat(value.([]interface{})[2].(string), 64)
		}
		if data["t"] != last {
			return errors.New("Last timestamp return by request if different from parsing")
		}
		accumulator.AddFields("spread", data, map[string]string{"asset": pair})
	}
	return nil
}

func init() {
	inputs.Add("spread", func() telegraf.Input {
		return NewSpread()
	})
}
