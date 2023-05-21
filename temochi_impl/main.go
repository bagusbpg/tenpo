package main

import (
	"os"

	"github.com/bagusbpg/tenpo/kikai/daemon"
	"github.com/bagusbpg/tenpo/temochi_impl/config"
	"github.com/bagusbpg/tenpo/temochi_impl/server"
	"golang.org/x/exp/slog"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)))
}

func main() {
	runner := daemon.NewServiceRunner()
	runner.Run(&server.Component{}, &config.Config{})
}
