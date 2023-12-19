package client

import (
	"net/http"

	tenpo_http "github.com/bagusbpg/tenpo/kikai/http"
)

type client struct {
	httpClient *http.Client
	url        string
}

type Config struct {
	ClientConfig tenpo_http.ClientConfig
	URL          string
}

func NewTemochiClient(config Config) *client {
	return &client{
		httpClient: tenpo_http.NewHTTPClient(config.ClientConfig),
		url:        config.URL,
	}
}
