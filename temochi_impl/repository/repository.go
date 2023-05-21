package repository

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetStocks(ctx context.Context, input GetStocksDBInput, output *GetStocksDBOutput) error
	UpsertStock(ctx context.Context, input UpsertStockDBInput, output *UpsertStockDBOutput) error
	UpdateChannelStocks(ctx context.Context, input UpdateChannelStocksDBInput, output *UpdateChannelStocksDBOutput) error
	DeleteChannelStock(ctx context.Context, input DeleteChannelStockDBInput, output *DeleteChannelStockDBOutput) error
	DeleteStock(ctx context.Context, input DeleteStockDBInput, output *DeleteStockDBOutput) error
}

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
