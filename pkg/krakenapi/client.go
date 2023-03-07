package krakenapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const UrlBase = "https://api.kraken.com/0/"

var Headers = map[string]string{"User-Agent": "telegraf-kraken"}

var ErrUnauthorized = errors.New("Missing or invalid API key")

type Client struct {
	Method    string
	UrlSuffix string
	APIKey    string

	c *http.Client
}

type Response struct {
	Result map[string]interface{} `json:"result"`
	Error  []string               `json:"error"`
}

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

func (c *Client) Request(headers map[string]string, parameters map[string]string) (map[string]interface{}, error) {
	req, err := http.NewRequest(c.Method, joinUrl(UrlBase, c.UrlSuffix), nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()
	// parsing data
	result := &Response{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if len(result.Error) > 0 {
		return nil, errors.New(strings.Join(result.Error, ","))
	}
	return result.Result, err
}

func joinUrl(base string, ext string) string {
	urlReq, _ := url.Parse(base + ext)
	return urlReq.String()
}
