package handler

import (
	"net/http"

	"github.com/bagusbpg/tenpo/temochi"
	"github.com/go-playground/validator/v10"
)

type Handler interface {
	GetStocks() http.HandlerFunc
	UpsertStocks() http.HandlerFunc
	UpdateChannelStocks() http.HandlerFunc
	DeleteChannelStock() http.HandlerFunc
	DeleteStock() http.HandlerFunc
}

type handler struct {
	service   temochi.Service
	validator *validator.Validate
}

func New(service temochi.Service, validator *validator.Validate) Handler {
	return &handler{
		service:   service,
		validator: validator,
	}
}
