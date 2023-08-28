package handler

import (
	"github.com/bagusbpg/tenpo/temochi"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	service   temochi.Service
	validator *validator.Validate
}

func New(service temochi.Service, validator *validator.Validate) *handler {
	return &handler{
		service:   service,
		validator: validator,
	}
}
