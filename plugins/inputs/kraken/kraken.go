//go:generate ../../../tools/readme_config_includer/generator
package kraken

import (
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"
	jsonparser "github.com/influxdata/telegraf/plugins/parsers/json"
)

//go:embed sample.conf
var sampleConfig string

const suffixTicker = "/Ticker"

type Ticker struct {
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

type Data struct {
	Result map[string]Ticker `json:"result"`
	Error  []string          `json:"error"`
}

type Kraken struct {
	URL      string            `toml:"url"`
	Pairs    []string          `toml:"pairs"`
	Includes []string          `toml:"include"`
	Method   string            `toml:"method"`
	Headers  map[string]string `toml:"headers"`
	Timeout  config.Duration   `toml:"timeout"`

	client *http.Client
}

func NewKraken() *Kraken {
	return &Kraken{
		URL:     "https://api.kraken.com/0/public",
		Pairs:   []string{"XRPUSDT", "ETHUSDC"},
		Method:  "GET",
		Headers: map[string]string{"User-Agent": "telegraf-kraken"},
		Timeout: config.Duration(time.Second * 5),
	}
}

func (*Kraken) SampleConfig() string {
	return sampleConfig
}

func (*Kraken) Description() string {
	return "Example go-plugin for Telegraf"
}

func (k *Kraken) Init() error {
	var err error
	k.client, err = k.createHTTPClient()
	if err != nil {
		return err
	}
	return nil
}

func (*Kraken) createHTTPClient() (*http.Client, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
		Timeout: time.Duration(cryptodl.Timeout),
	}
	return client, nil
}

// gatherJSONData query the data source and parse the response JSON
func (k *Kraken) gatherJSONData(address string, parameters map[string]string, value interface{}) error {
	request, err := http.NewRequest(k.Method, address, nil)
	if err != nil {
		return err
	}
	// headers
	for key, values := range k.Headers {
		request.Header.Add(key, values)
	}
	// parameters
	query := request.URL.Query()
	for label, parameter := range parameters {
		query.Add(label, parameter)
	}
	request.URL.RawQuery = query.Encode()
	// request
	response, err := k.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(value)
}

func (k *Kraken) Gather(accumulator telegraf.Accumulator) error {
	data := &Data{}
	// url
	tickerURL, err := url.Parse(cryptodl.URL + suffixTicker)
	if err != nil {
		return err
	}
	// parameter
	parameters := map[string]string{"pair": strings.Join(k.Pairs, ",")}
	// request
	err = cryptodl.gatherJSONData(tickerURL.String(), parameters, data)
	if err != nil {
		return err
	}
	if len(data.Error) > 0 {
		return errors.New(strings.Join(data.Error, ","))
	}
	// aggregate
	for pair := range data.Result {
		var result map[string]interface{}
		jrec, _ := json.Marshal(data.Result[pair])
		json.Unmarshal(jrec, &record)

		flattener := jsonparser.JSONFlattener{}
		err := flattener.FullFlattenJSON("", record, true, true)
		if err != nil {
			return err
		}
		accumulator.AddFields("ticker", flattener.Fields, map[string]string{"asset": pair})
	}
	return nil
}

func init() {
	inputs.Add("kraken", func() telegraf.Input {
		return NewKraken
	})
}
