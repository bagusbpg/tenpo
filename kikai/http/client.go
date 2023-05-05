package http

import (
	"net/http"
	"time"
)

type ClientConfig struct {
	TimeoutSec          int
	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
}

func NewHTTPClient(config ClientConfig) *http.Client {
	requestTimeout := DEFAULT_HTTP_CLIENT_REQUEST_TIMEOUT
	if config.TimeoutSec != 0 {
		requestTimeout = time.Duration(config.TimeoutSec) * time.Second
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if config.MaxIdleConns != 0 {
		transport.MaxIdleConns = config.MaxIdleConns
	}
	if config.MaxConnsPerHost != 0 {
		transport.MaxConnsPerHost = config.MaxConnsPerHost
	}
	if config.MaxIdleConnsPerHost != 0 {
		transport.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost
	}

	return &http.Client{
		Timeout:   requestTimeout,
		Transport: transport,
	}
}
