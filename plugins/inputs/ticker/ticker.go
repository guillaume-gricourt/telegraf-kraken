//go:generate ../../../tools/readme_config_includer/generator
package ticker

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

const suffixQuery = "/Ticker"

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
	URL      string            `toml:"url"`
	Pairs    []string          `toml:"pairs"`
	Includes []string          `toml:"include"`
	Method   string            `toml:"method"`
	Headers  map[string]string `toml:"headers"`
	Timeout  config.Duration   `toml:"timeout"`

	client *http.Client
}

func NewTicker() *Ticker {
	return &Ticker{
		URL:     "https://api.kraken.com/0/public",
		Pairs:   []string{},
		Method:  "GET",
		Headers: map[string]string{"User-Agent": "telegraf-kraken"},
		Timeout: config.Duration(time.Second * 5),
	}
}

func (*Ticker) SampleConfig() string {
	return sampleConfig
}

func (*Ticker) Description() string {
	return "Request for Kraken - Ticker"
}

func (t *Ticker) Init() error {
	var err error
	t.client, err = t.createHTTPClient()
	if err != nil {
		return err
	}
	if len(t.Pairs) < 1 {
		return errors.New("Provide at least one asset to download")
	}
	return nil
}

func (t *Ticker) createHTTPClient() (*http.Client, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
		Timeout: time.Duration(t.Timeout),
	}
	return client, nil
}

// gatherJSONData query the data source and parse the response JSON
func (t *Ticker) gatherJSONData(address string, parameters map[string]string, value interface{}) error {
	request, err := http.NewRequest(t.Method, address, nil)
	if err != nil {
		return err
	}
	// headers
	for key, values := range t.Headers {
		request.Header.Add(key, values)
	}
	// parameters
	query := request.URL.Query()
	for label, parameter := range parameters {
		query.Add(label, parameter)
	}
	request.URL.RawQuery = query.Encode()
	// request
	response, err := t.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(value)
}

func (t *Ticker) Gather(accumulator telegraf.Accumulator) error {
	resp := &Response{}
	// url
	tickerURL, err := url.Parse(t.URL + suffixQuery)
	if err != nil {
		return err
	}
	// parameter
	parameters := map[string]string{"pair": strings.Join(t.Pairs, ",")}
	// request
	err = t.gatherJSONData(tickerURL.String(), parameters, resp)
	if err != nil {
		return err
	}
	if len(resp.Error) > 0 {
		return errors.New(strings.Join(resp.Error, ","))
	}
	// aggregate
	for pair := range resp.Result {
		var record map[string]interface{}
		jrec, _ := json.Marshal(resp.Result[pair])
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
	inputs.Add("ticker", func() telegraf.Input {
		return NewTicker()
	})
}
