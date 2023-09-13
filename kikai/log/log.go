package log

import (
	"context"
	"log/slog"
	"net/http"

	tenpoHttp "github.com/bagusbpg/tenpo/kikai/http"
)

func Error(r *http.Request, errMessage string, cause error) {
	slog.LogAttrs(
		r.Context(),
		slog.LevelError, errMessage,
		slog.String("causedBy", cause.Error()),
		slog.String("path", r.URL.Path),
		slog.String("requestID", getRequestID(r.Context())),
	)
}

func Success(r *http.Request, successInfo string, handler string) {
	slog.LogAttrs(
		r.Context(),
		slog.LevelInfo, successInfo,
		slog.String("handler", handler),
		slog.String("path", r.URL.Path),
		slog.String("requestID", getRequestID(r.Context())),
	)
}

func getRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(tenpoHttp.REQUEST_ID_CONTEXT_KEY).(string)

	return requestID
}
