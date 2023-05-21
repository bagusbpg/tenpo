package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

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

func (ths *repository) UpsertStock(ctx context.Context, input UpsertStockDBInput, output *UpsertStockDBOutput) error {
	tx, err := ths.db.Begin()
	if err != nil {
		return fmt.Errorf("failed starting UpsertStock transaction: %s", err.Error())
	}
	defer tx.Rollback()

	queryUpsertInventory, argsUpsertInventory := buildUpsertInventoryQuery(input)
	_, err = tx.ExecContext(ctx, queryUpsertInventory, argsUpsertInventory...)
	if err != nil {
		return fmt.Errorf("failed executing UpsertInventory query: %s", err.Error())
	}

	if len(input.UpsertChannelStockInputs) > 0 {
		queryUpsertChannelStock, argsUpsertChannelStock := buildUpsertChannelStockQuery(input)
		_, err = tx.ExecContext(ctx, queryUpsertChannelStock, argsUpsertChannelStock...)
		if err != nil {
			return fmt.Errorf("failed executing UpsertChannelStock query: %s", err.Error())
		}

		queryDeleteChannelStock, argsDeleteChannelStock := buildDeleteExcludedGateChannelStockQuery(input)
		_, err = tx.ExecContext(ctx, queryDeleteChannelStock, argsDeleteChannelStock...)
		if err != nil {
			return fmt.Errorf("failed executing DeletedChannelStock query: %s", err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed committing UpsertStock transaction: %s", err.Error())
	}

	return nil
}

func buildUpsertInventoryQuery(input UpsertStockDBInput) (string, []interface{}) {
	queryBuilder := sq.
		Insert(`"temochi".inventory`).
		Columns(
			"warehouse_id",
			"sku",
			"stock",
			"buffer_stock",
		).
		Suffix("ON CONFLICT ON CONSTRAINT inventory_pk DO UPDATE SET stock = EXCLUDED.stock, buffer_stock = EXCLUDED.buffer_stock, version = inventory.version + 1, updated_at = NOW()")

	for _, item := range input.UpsertInventoryInputs {
		queryBuilder = queryBuilder.Values(input.WarehouseID, item.SKU, item.Stock, item.BufferStock)
	}

	return queryBuilder.PlaceholderFormat(sq.Dollar).MustSql()
}

func buildUpsertChannelStockQuery(input UpsertStockDBInput) (string, []interface{}) {
	queryBuilder := sq.
		Insert(`"temochi".channel_stock`).
		Columns(
			"warehouse_id",
			"sku",
			"gate_id",
			"channel_id",
			"stock",
		).
		Suffix("ON CONFLICT ON CONSTRAINT channel_stock_pk DO UPDATE SET stock = EXCLUDED.stock, version = channel_stock.version + 1, updated_at = NOW()")

	for _, item := range input.UpsertChannelStockInputs {
		queryBuilder = queryBuilder.Values(input.WarehouseID, item.SKU, item.GateID, item.ChannelID, item.Stock)
	}

	return queryBuilder.PlaceholderFormat(sq.Dollar).MustSql()
}

func buildDeleteExcludedGateChannelStockQuery(input UpsertStockDBInput) (string, []interface{}) {
	excludedGateChannel := make([]string, 0, len(input.UpsertChannelStockInputs))
	skus := make([]string, 0, len(input.UpsertChannelStockInputs))
	for _, item := range input.UpsertChannelStockInputs {
		excludedGateChannel = append(excludedGateChannel, item.SKU+"#"+item.GateID+"#"+item.ChannelID)
		skus = append(skus, item.SKU)
	}

	queryBuilder := sq.
		Delete(`"temochi".channel_stock`).
		Where(sq.And{
			sq.Eq{"warehouse_id": input.WarehouseID},
			sq.Eq{"sku": skus},
			sq.NotEq{"CONCAT(sku, '#', gate_id, '#', channel_id)": excludedGateChannel},
		})

	return queryBuilder.PlaceholderFormat(sq.Dollar).MustSql()
}
