package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var testBaseURL = url.URL{
	Scheme: "https",
	Host:   "testapi.payarc.net",
	Path:   "/v1",
}

var baseURL = url.URL{
	Scheme: "https",
	Host:   "api.payarc.net",
	Path:   "/v1",
}

type Client struct {
	ApiKey     string
	HttpClient *http.Client
	Url        *url.URL
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("no api key provided")
	}

	client := &Client{
		ApiKey: apiKey,
		HttpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	if true {
		client.Url = &testBaseURL
	} else {
		client.Url = &baseURL
	}

	return client, nil
}
