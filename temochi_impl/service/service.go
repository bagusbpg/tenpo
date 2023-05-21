package service

import (
	"github.com/bagusbpg/tenpo/temochi"
	"github.com/bagusbpg/tenpo/temochi_impl/repository"
)

type service struct {
	repository repository.Repository
}

func New(repository repository.Repository) temochi.Service {
	return &service{
		repository: repository,
	}
}
