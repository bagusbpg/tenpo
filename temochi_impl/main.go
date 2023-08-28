package main

import (
	"log/slog"
	"os"

	"github.com/bagusbpg/tenpo/kikai/daemon"
	"github.com/bagusbpg/tenpo/temochi_impl/config"
	"github.com/bagusbpg/tenpo/temochi_impl/server"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

func main() {
	daemon.Run(&server.Component{}, &config.Config{})
}
