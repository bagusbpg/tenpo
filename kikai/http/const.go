package http

import "time"

const (
	DEFAULT_HTTP_SERVER_PORT                       string        = "8000"
	DEFAULT_HTTP_SERVER_GRACEFULL_SHUTDOWN_TIMEOUT time.Duration = 10
	DEFAULT_HTTP_CLIENT_REQUEST_TIMEOUT            time.Duration = 0
)
