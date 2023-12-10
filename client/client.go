package client

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

type PayArcEnvironment int

const (
	PayArcEnvironmentTest       PayArcEnvironment = iota
	PayArcEnvironmentProduction PayArcEnvironment = iota
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

func NewClient(apiKey string, environment PayArcEnvironment) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("no api key provided")
	}

	client := &Client{
		ApiKey: apiKey,
		HttpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	switch environment {
	case PayArcEnvironmentTest:
		client.Url = &testBaseURL
	case PayArcEnvironmentProduction:
		client.Url = &baseURL
	}

	return client, nil
}
