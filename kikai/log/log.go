package log

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	tenpoHttp "github.com/bagusbpg/tenpo/kikai/http"
)

func Error(r *http.Request, handler string, cause error) {
	slog.LogAttrs(
		r.Context(),
		slog.LevelError, fmt.Sprintf("failed processing request at %s", handler),
		slog.String("causedBy", cause.Error()),
		slog.String("path", r.URL.Path),
		slog.String("requestID", getRequestID(r.Context())),
	)
}

func Success(r *http.Request, handler string) {
	slog.LogAttrs(
		r.Context(),
		slog.LevelInfo, fmt.Sprintf("success processing request at %s", handler),
		slog.String("path", r.URL.Path),
		slog.String("requestID", getRequestID(r.Context())),
	)
}

func getRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(tenpoHttp.REQUEST_ID_CONTEXT_KEY).(string)

	return requestID
}
