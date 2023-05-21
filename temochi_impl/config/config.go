package config

import (
	tenpoHttp "github.com/bagusbpg/tenpo/kikai/http"
	tenpoSql "github.com/bagusbpg/tenpo/kikai/sql"
)

type Config struct {
	ServerConfig tenpoHttp.ServerConfig
	DBConfig     tenpoSql.Config
}
