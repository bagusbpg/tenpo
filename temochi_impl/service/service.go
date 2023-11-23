package service

import (
	"context"

	"github.com/bagusbpg/tenpo/temochi_impl/repository"
)

type service struct {
	repository Repository
}

//go:generate mockgen -source=./service.go -destination=./mock.go -package=service
type Repository interface {
	GetStocks(ctx context.Context, input repository.GetStocksDBInput, output *repository.GetStocksDBOutput) error
	UpsertStock(ctx context.Context, input repository.UpsertStockDBInput, output *repository.UpsertStockDBOutput) error
	UpdateChannelStocks(ctx context.Context, input repository.UpdateChannelStocksDBInput, output *repository.UpdateChannelStocksDBOutput) error
	DeleteChannelStock(ctx context.Context, input repository.DeleteChannelStockDBInput, output *repository.DeleteChannelStockDBOutput) error
	DeleteStock(ctx context.Context, input repository.DeleteStockDBInput, output *repository.DeleteStockDBOutput) error
}

func New(repository Repository) *service {
	return &service{
		repository: repository,
	}
}
