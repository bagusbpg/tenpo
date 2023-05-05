package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/slog"
)

type Component interface {
	// New initiates the service with the given config
	New(config interface{}) error

	// Start starts the service's components
	Start() error

	// Stop stops the service's components gracefully
	Stop() error
}

type Runner interface {
	// Run will call New and then Start. Run will call
	// Stop when syscall.SIGINT or syscall.SIGTERM is
	// received.
	Run(component Component, config interface{})
}

type runner struct{}

func NewServiceRunner() Runner {
	return &runner{}
}

func (ths *runner) Run(component Component, config interface{}) {
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminated := make(chan struct{}, 1)
	go func() {
		sig := <-c
		slog.LogAttrs(
			ctx,
			slog.LevelInfo, "exiting service",
			slog.String("receivingSignal", sig.String()),
		)

		if err := component.Stop(); err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError, "failed stopping service",
				slog.String("causedBy", err.Error()),
			)
		}

		close(terminated)
	}()

	envMode := os.Getenv("ENV_MODE")
	if !isValidEnvMode(envMode) {
		slog.LogAttrs(
			ctx,
			slog.LevelError, "invalid environment mode",
			slog.String("ENV_MODE", envMode),
		)
		return
	}

	err := loadConfig(&config, envMode)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError, "failed loading config",
			slog.String("causedBy", err.Error()),
		)

		return
	}

	err = component.New(config)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError, "failed initiating service",
			slog.String("causedBy", err.Error()),
		)

		return
	}

	slog.LogAttrs(
		ctx,
		slog.LevelInfo, "staring service",
		slog.String("ENV_MODE", envMode),
	)

	err = component.Start()
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError, "failed starting service",
			slog.String("causedBy", err.Error()),
		)

		err = component.Stop()
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError, "failed stopping service",
				slog.String("causedBy", err.Error()),
			)
		}

		close(terminated)
	}

	<-terminated
	slog.LogAttrs(
		ctx,
		slog.LevelInfo, "service stopped gracefully",
	)
}

func isValidEnvMode(envMode string) bool {
	switch envMode {
	case "dev", "stg", "prod":
		return true
	default:
		return false
	}
}

func loadConfig(config interface{}, envMode string) error {
	file := fmt.Sprintf("./config/%s.json", envMode)
	jsonFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error opening file [%s], cause: %s", file, err)
	}
	defer jsonFile.Close()

	jsonByte, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("error reading file [%s], cause: %s", file, err)
	}

	err = json.Unmarshal(jsonByte, config)
	if err != nil {
		return fmt.Errorf("error unmarshal json, cause: %s", err)
	}

	return nil
}
