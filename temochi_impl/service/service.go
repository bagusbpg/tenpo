package service

import (
	"context"

	"github.com/bagusbpg/tenpo/temochi"
)

type service struct {
	repository Repository
}

type Repository interface {
	GetStocks(ctx context.Context, input GetStocksDBInput, output *GetStocksDBOutput) error
	UpsertStock(ctx context.Context, input UpsertStockDBInput, output *UpsertStockDBOutput) error
	UpdateChannelStocks(ctx context.Context, input UpdateChannelStocksDBInput, output *UpdateChannelStocksDBOutput) error
	DeleteChannelStock(ctx context.Context, input DeleteChannelStockDBInput, output *DeleteChannelStockDBOutput) error
	DeleteStock(ctx context.Context, input DeleteStockDBInput, output *DeleteStockDBOutput) error
}

type GetStocksDBInput struct {
	WarehouseID string
	SKUs        []string
}

type GetStocksDBOutput struct {
	Stocks []StockDB
}

type StockDB struct {
	temochi.Inventory
	ChannelStocks []temochi.ChannelStock `json:"channelStocks"`
}

type UpsertStockDBInput struct {
	WarehouseID              string
	UpsertInventoryInputs    []UpsertInventoryInput
	UpsertChannelStockInputs []UpsertChannelStockInput
}

type UpsertInventoryInput struct {
	SKU         string
	Stock       uint32
	BufferStock uint32
}

type UpsertChannelStockInput struct {
	SKU       string
	GateID    string
	ChannelID string
	Stock     uint32
}

type UpsertStockDBOutput struct{}

type UpdateChannelStocksDBInput struct {
	WarehouseID              string
	UpdateChannelStockInputs []UpdateChannelStockInput
}

type UpdateChannelStockInput struct {
	SKU       string
	GateID    string
	ChannelID string
	Delta     int32
}

type UpdateChannelStocksDBOutput struct{}

type DeleteChannelStockDBInput struct {
	WarehouseID string
	GateID      string
	ChannelID   string
}

type DeleteChannelStockDBOutput struct{}

type DeleteStockDBInput struct {
	WarehouseID string
	SKU         string
}

type DeleteStockDBOutput struct{}

func New(repository Repository) *service {
	return &service{
		repository: repository,
	}
}
