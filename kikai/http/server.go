package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/exp/slog"
)

type Server interface {
	// Start starts http server by calling ListenAndServe
	Start() error

	// Stop stops http server by calling Shutdown
	Stop() error

	// AddRoute registers new route to http router with
	// default logger and response writer content-type
	// (json) are set via global middleware
	AddRoute(method string, pattern string, handler http.Handler)
}

type server struct {
	httpServer      *http.Server
	router          *httprouter.Router
	shutdownTimeout time.Duration
}

type ServerConfig struct {
	GracefulShutdownTimeoutSec uint8
}

func NewHTTPServer(config ServerConfig) Server {
	ctx := context.Background()
	port := DEFAULT_HTTP_SERVER_PORT
	httpServerPort := os.Getenv("HTTP_SERVER_PORT")
	if httpServerPort != "" && httpServerPort != DEFAULT_HTTP_SERVER_PORT {
		slog.LogAttrs(
			ctx,
			slog.LevelInfo, "HTTP_SERVER_PORT environment variable is detected",
			slog.String("HTTP_SERVER_PORT", httpServerPort),
		)

		if regexp.MustCompile(`^[1-9][0-9]*$`).MatchString(httpServerPort) {
			slog.LogAttrs(
				ctx,
				slog.LevelWarn, "invalid HTTP_SERVER_PORT environment variable, fallback to DEFAULT_HTTP_SERVER_PORT",
			)
		} else {
			port = httpServerPort
		}
	}

	shutdownTimeout := DEFAULT_HTTP_SERVER_GRACEFULL_SHUTDOWN_TIMEOUT * time.Second
	if config.GracefulShutdownTimeoutSec != 0 {
		shutdownTimeout = time.Duration(config.GracefulShutdownTimeoutSec) * time.Second
	}

	router := httprouter.New()

	return &server{
		httpServer: &http.Server{
			Addr:    "0.0.0.0:" + port,
			Handler: router,
		},
		router:          router,
		shutdownTimeout: shutdownTimeout,
	}
}

func (ths *server) Start() error {
	ctx := context.Background()
	slog.LogAttrs(
		ctx,
		slog.LevelInfo, "starting http server",
		slog.String("address", ths.httpServer.Addr),
	)

	if err := ths.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed listening and serving http server: %s", err.Error())
	}

	return nil
}

func (ths *server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), ths.shutdownTimeout)
	defer cancel()

	if err := ths.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed shutting down http server: %s", err.Error())
	}

	slog.LogAttrs(
		ctx,
		slog.LevelInfo, "http server stopped gracefully",
	)

	return nil
}

func (ths *server) AddRoute(method string, path string, handler http.Handler) {
	handlerName := getFuncName(handler)
	newHandler := func(h http.Handler) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			w.Header().Set("Content-Type", "application/json")

			ctx := r.Context()
			requestID, ok := ctx.Value(REQUEST_ID_CONTEXT_KEY).(string)
			if !ok {
				requestID = uuid.New().String()
				ctx = context.WithValue(ctx, REQUEST_ID_CONTEXT_KEY, requestID)
				r = r.WithContext(ctx)
			}
			slog.LogAttrs(
				context.Background(),
				slog.LevelInfo, "receiving http request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.String()),
				slog.String("handler", handlerName),
				slog.String("requestID", requestID),
			)

			handler.ServeHTTP(w, r)
		}
	}(handler)

	switch method {
	case http.MethodPost:
		ths.router.POST(path, newHandler)
	case http.MethodGet:
		ths.router.GET(path, newHandler)
	case http.MethodPut:
		ths.router.PUT(path, newHandler)
	case http.MethodPatch:
		ths.router.PATCH(path, newHandler)
	case http.MethodDelete:
		ths.router.DELETE(path, newHandler)
	default:
		slog.LogAttrs(
			context.Background(),
			slog.LevelWarn, "ignoring unsupported http method",
			slog.String("method", method),
		)
	}
}

var re = regexp.MustCompile(`\)\.[^\.]*`)

func getFuncName(handler http.Handler) string {
	return strings.TrimPrefix(re.FindString(runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()), ").")
}
