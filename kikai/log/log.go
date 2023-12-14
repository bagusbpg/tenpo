package log

import (
	"context"
	"log/slog"

	tenpoHttp "github.com/bagusbpg/tenpo/kikai/http"
)

func Error(ctx context.Context, path string, err error) {
	slog.LogAttrs(
		ctx,
		slog.LevelError, err.Error(),
		slog.String("path", path),
		slog.String("requestID", getRequestID(ctx)),
	)
}

func Success(ctx context.Context, path string) {
	slog.LogAttrs(
		ctx,
		slog.LevelInfo, "success",
		slog.String("path", path),
		slog.String("requestID", getRequestID(ctx)),
	)
}

func getRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(tenpoHttp.REQUEST_ID_CONTEXT_KEY).(string)

	return requestID
}
