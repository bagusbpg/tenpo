package server

import (
	"fmt"
	"net/http"

	tenpoHttp "github.com/bagusbpg/tenpo/kikai/http"
	tenpoSql "github.com/bagusbpg/tenpo/kikai/sql"
	"github.com/bagusbpg/tenpo/temochi"
	"github.com/bagusbpg/tenpo/temochi_impl/config"
	"github.com/bagusbpg/tenpo/temochi_impl/handler"
	"github.com/bagusbpg/tenpo/temochi_impl/repository"
	"github.com/bagusbpg/tenpo/temochi_impl/service"
	"github.com/go-playground/validator/v10"
)

type Component struct {
	config *config.Config
	server *tenpoHttp.Server
}

func (ths *Component) New(appConfig interface{}) error {
	ths.config = appConfig.(*config.Config)
	ths.server = tenpoHttp.NewHTTPServer(ths.config.ServerConfig)

	db, err := tenpoSql.NewClient(ths.config.DBConfig)
	if err != nil {
		return fmt.Errorf("failed opening database connection: %s", err.Error())
	}

	repository := repository.New(db)
	service := service.New(repository)
	validator := validator.New()
	handler := handler.New(service, validator)

	ths.server.AddRoute(http.MethodGet, temochi.PATH_GET_STOCKS, "GetStocks", handler.GetStocks())
	ths.server.AddRoute(http.MethodPost, temochi.PATH_UPSERT_STOCKS, "UpsertStocks", handler.UpsertStocks())
	ths.server.AddRoute(http.MethodPut, temochi.PATH_UPDATE_CHANNELS_STOCK, "UpdateChannelStocks", handler.UpdateChannelStocks())
	ths.server.AddRoute(http.MethodPatch, temochi.PATH_DELETE_CHANNEL_STOCK, "DeleteChannelStock", handler.DeleteChannelStock())
	ths.server.AddRoute(http.MethodDelete, temochi.PATH_DELETE_STOCK, "DeleteStock", handler.DeleteStock())

	return nil
}

func (ths *Component) Start() error {
	if err := ths.server.Start(); err != nil {
		return fmt.Errorf("failed starting http server: %s", err.Error())
	}

	return nil
}

func (ths *Component) Stop() error {
	if err := ths.server.Stop(); err != nil {
		return fmt.Errorf("failed stopping http server: %s", err.Error())
	}

	return nil
}
