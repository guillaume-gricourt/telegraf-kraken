package krakenapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// Values
const UrlBase = "https://api.kraken.com/0/"

var Headers = map[string]string{"User-Agent": "telegraf-kraken"}

var ErrUnauthorized = errors.New("Missing or invalid API key")

// A Client manages communication with the OctoPrint API.
type Client struct {
	Method    string
	UrlSuffix string
	APIKey    string

	c *http.Client
}

// NewClient returns a new OctoPrint API client with provided base URL and API
// Key. If baseURL does not have a trailing slash, one is added automatically. If
// `Access Control` is enabled at OctoPrint configuration an apiKey should be
// provided (http://docs.octoprint.org/en/master/api/general.html#authorization).
func NewClient(method string, urlSuffix string, apiKey string, timeOut time.Duration) *Client {
	return &Client{
		Method:    method,
		UrlSuffix: urlSuffix,
		APIKey:    apiKey,
		c: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			Timeout: time.Duration(timeOut),
		},
	}
}

// gatherJSONData query the data source and parse the response JSON
func (c *Client) Request(headers map[string]string, parameters map[string]string, value interface{}) error {
	req, err := http.NewRequest(c.Method, joinUrl(UrlBase, c.UrlSuffix), nil)
	if err != nil {
		return err
	}
	// headers
	if headers != nil {
		for key, values := range headers {
			Headers[key] = values
		}
	}
	for key, values := range Headers {
		req.Header.Add(key, values)
	}
	// parameters
	if parameters != nil {
		query := req.URL.Query()
		for label, parameter := range parameters {
			query.Add(label, parameter)
		}
		req.URL.RawQuery = query.Encode()
	}
	// request
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(value)
}

func joinUrl(base string, ext string) string {
	urlReq, _ := url.Parse(base + ext)
	return urlReq.String()
}
