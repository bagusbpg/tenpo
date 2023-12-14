package log

import (
	"context"
	"log/slog"

	tenpoHttp "github.com/bagusbpg/tenpo/kikai/http"
)

func Warn(ctx context.Context, message string) {
	slog.LogAttrs(
		ctx,
		slog.LevelWarn, message,
		slog.String("requestID", getRequestID(ctx)),
		slog.String("path", getRequestPath(ctx)),
		slog.String("method", getRequestMethod(ctx)),
		slog.String("handler", getRequestHandler(ctx)),
	)
}

func Error(ctx context.Context, err error) {
	slog.LogAttrs(
		ctx,
		slog.LevelError, err.Error(),
		slog.String("requestID", getRequestID(ctx)),
		slog.String("path", getRequestPath(ctx)),
		slog.String("method", getRequestMethod(ctx)),
		slog.String("handler", getRequestHandler(ctx)),
	)
}

func Success(ctx context.Context) {
	slog.LogAttrs(
		ctx,
		slog.LevelInfo, "success",
		slog.String("requestID", getRequestID(ctx)),
		slog.String("path", getRequestPath(ctx)),
		slog.String("method", getRequestMethod(ctx)),
		slog.String("handler", getRequestHandler(ctx)),
	)
}

func getRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(tenpoHttp.REQUEST_ID_CONTEXT_KEY).(string)

	return requestID
}

func getRequestPath(ctx context.Context) string {
	path, _ := ctx.Value(tenpoHttp.REQUEST_PATH_CONTEXT_KEY).(string)

	return path
}

func getRequestMethod(ctx context.Context) string {
	method, _ := ctx.Value(tenpoHttp.REQUEST_METHOD_CONTEXT_KEY).(string)

	return method
}

func getRequestHandler(ctx context.Context) string {
	handler, _ := ctx.Value(tenpoHttp.REQUEST_HANDLER_CONTEXT_KEY).(string)

	return handler
}
