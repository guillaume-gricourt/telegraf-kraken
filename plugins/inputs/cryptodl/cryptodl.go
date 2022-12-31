//go:generate ../../../tools/readme_config_includer/generator
package cryptodl

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/internal/choice"
	"github.com/influxdata/telegraf/plugins/common/tls"
	"github.com/influxdata/telegraf/plugins/inputs"
	jsonparser "github.com/influxdata/telegraf/plugins/parsers/json"
)

//go:embed sample.conf
var sampleConfig string

const suffixInfo = "/"
const suffixStats = "/stats"

type Info struct {
	CryptoDl string `json:"cryptodl"`
	Hostname string `json:"hostname"`
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	Version  string `json:"version"`
}

type Stats struct {
	CryptoDl     map[string]interface{} `json:"cryptodl"`
	FileCryptoDl interface{}            `json:"filecryptodl"`
	Libcryptodl  interface{}            `json:"libcryptodl"`
	System       interface{}            `json:"system"`
}

type CryptoDl struct {
	URL string `toml:"url"`

	Includes []string `toml:"include"`

	Username   string            `toml:"username"`
	Password   string            `toml:"password"`
	Method     string            `toml:"method"`
	Headers    map[string]string `toml:"headers"`
	HostHeader string            `toml:"host_header"`
	Timeout    config.Duration   `toml:"timeout"`

	tls.ClientConfig
	client *http.Client
}

func NewCryptoDl() *CryptoDl {
	return &CryptoDl{
		URL:      "http://127.0.0.1",
		Includes: []string{"cryptodl", "libcryptodl", "filecryptodl"},
		Method:   "GET",
		Headers:  make(map[string]string),
		Timeout:  config.Duration(time.Second * 5),
	}
}

func (*CryptoDl) SampleConfig() string {
	return sampleConfig
}

// Init is for setup, and validating config.
func (s *CryptoDl) Init() error {
	availableStats := []string{"cryptodl", "libcryptodl", "system", "filecryptodl"}

	var err error
	cryptodl.client, err = cryptodl.createHTTPClient()

	if err != nil {
		return err
	}

	err = choice.CheckSlice(cryptodl.Includes, availableStats)
	if err != nil {
		return err
	}

	return nil
}

// createHTTPClient create a clients to access API
func (cryptodl *CryptoDl) createHTTPClient() (*http.Client, error) {
	tlsConfig, err := cryptodl.ClientConfig.TLSConfig()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: time.Duration(cryptodl.Timeout),
	}

	return client, nil
}

// gatherJSONData query the data source and parse the response JSON
func (cryptodl *CryptoDl) gatherJSONData(address string, value interface{}) error {
	request, err := http.NewRequest(cryptodl.Method, address, nil)
	if err != nil {
		return err
	}

	if cryptodl.Username != "" {
		request.SetBasicAuth(cryptodl.Username, cryptodl.Password)
	}
	for k, v := range cryptodl.Headers {
		request.Header.Add(k, v)
	}
	if cryptodl.HostHeader != "" {
		request.Host = cryptodl.HostHeader
	}

	response, err := cryptodl.client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(value)
}

func (cryptodl *Beat) Gather(accumulator telegraf.Accumulator) error {
	cryptodlStats := &Stats{}
	cryptodlInfo := &Info{}

	infoURL, err := url.Parse(cryptodl.URL + suffixInfo)
	if err != nil {
		return err
	}
	statsURL, err := url.Parse(cryptodl.URL + suffixStats)
	if err != nil {
		return err
	}

	err = cryptodl.gatherJSONData(infoURL.String(), cryptodlInfo)
	if err != nil {
		return err
	}
	tags := map[string]string{
		"cryptodl_cryptodl": cryptodlInfo.Beat,
		"cryptodl_id":       cryptodlInfo.UUID,
		"cryptodl_name":     cryptodlInfo.Name,
		"cryptodl_host":     cryptodlInfo.Hostname,
		"cryptodl_version":  cryptodlInfo.Version,
	}

	err = cryptodl.gatherJSONData(statsURL.String(), cryptodlStats)
	if err != nil {
		return err
	}

	for _, name := range cryptodl.Includes {
		var stats interface{}
		var metric string

		switch name {
		case "cryptodl":
			stats = cryptodlStats.Beat
			metric = "cryptodl"
		case "filecryptodl":
			stats = cryptodlStats.FileCryptoDl
			metric = "cryptodl_filecryptodl"
		case "system":
			stats = cryptodlStats.System
			metric = "cryptodl_system"
		case "libcryptodl":
			stats = cryptodlStats.Libcryptodl
			metric = "cryptodl_libcryptodl"
		default:
			return fmt.Errorf("unknown stats-type %q", name)
		}
		flattener := jsonparser.JSONFlattener{}
		err := flattener.FullFlattenJSON("", stats, true, true)
		if err != nil {
			return err
		}
		accumulator.AddFields(metric, flattener.Fields, tags)
	}

	return nil
}

func init() {
	inputs.Add("cryptodl", func() telegraf.Input {
		return NewCryptoDl()
	})
}
