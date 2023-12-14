package http

import "time"

type context_key string

const (
	DEFAULT_HTTP_SERVER_PORT                       string        = "8000"
	DEFAULT_HTTP_SERVER_GRACEFULL_SHUTDOWN_TIMEOUT time.Duration = 10
	DEFAULT_HTTP_CLIENT_REQUEST_TIMEOUT            time.Duration = 0

	REQUEST_ID_CONTEXT_KEY      context_key = "requestID"
	REQUEST_PATH_CONTEXT_KEY    context_key = "path"
	REQUEST_METHOD_CONTEXT_KEY  context_key = "method"
	REQUEST_HANDLER_CONTEXT_KEY context_key = "handler"
)
