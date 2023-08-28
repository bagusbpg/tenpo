package repository

import (
	"context"
	"fmt"

	"github.com/bagusbpg/tenpo/temochi_impl/service"
)

const UPDATE_CHANNEL_STOCK_QUERY = `
UPDATE "temochi".channel_stock
SET stock = channel_stock.stock + $1,
	version = channel_stock.version + 1,
	updated_at = NOW()
WHERE warehouse_id = $2 AND sku = $3 AND gate_id = $4 AND channel_id = $5`

const UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK = `
UPDATE "temochi".inventory
SET stock = inventory.stock + $1,
	buffer_stock = 
		CASE
			WHEN inventory.stock + $1 <= inventory.buffer_stock AND inventory.buffer_stock + $1 >= 0
				THEN inventory.buffer_stock + $1
			WHEN inventory.stock + $1 <= inventory.buffer_stock AND inventory.buffer_stock + $1 < 0
				THEN 0
			ELSE inventory.buffer_stock
		END,
	version = inventory.version + 1,
	updated_at = NOW()
WHERE warehouse_id = $2 and sku = $3`

const UPDATE_RELATED_CHANNEL_STOCK_QUERY = `
UPDATE "temochi".channel_stock
SET stock = channel_stock.stock + $1,
	version = channel_stock.version + 1,
	updated_at = NOW()
WHERE warehouse_id = $2 AND sku = $3 AND (gate_id <> $4 OR channel_id <> $5)`

func (ths *repository) UpdateChannelStocks(ctx context.Context, input service.UpdateChannelStocksDBInput, output *service.UpdateChannelStocksDBOutput) error {
	tx, err := ths.db.Begin()
	if err != nil {
		return fmt.Errorf("failed starting UpdateChannelStocks transaction: %s", err.Error())
	}
	defer tx.Rollback()

	stmtChannelStock, err := tx.PrepareContext(ctx, UPDATE_CHANNEL_STOCK_QUERY)
	if err != nil {
		return fmt.Errorf("failed preparing statement UpdateChannelStock: %s", err.Error())
	}

	stmtInventoryWithBufferStock, err := tx.PrepareContext(ctx, UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK)
	if err != nil {
		return fmt.Errorf("failed preparing statement UpdateInventory with buffer stock: %s", err.Error())
	}

	stmtRelatedChannelStock, err := tx.PrepareContext(ctx, UPDATE_RELATED_CHANNEL_STOCK_QUERY)
	if err != nil {
		return fmt.Errorf("failed preparing statement UpdateChannelStock for related stock: %s", err.Error())
	}

	for _, item := range input.UpdateChannelStockInputs {
		_, err = stmtChannelStock.Exec(item.Delta, input.WarehouseID, item.SKU, item.GateID, item.ChannelID)
		if err != nil {
			return fmt.Errorf("failed executing UpdateChannelStock query: %s", err.Error())
		}

		_, err = stmtInventoryWithBufferStock.Exec(item.Delta, input.WarehouseID, item.SKU)
		if err != nil {
			return fmt.Errorf("failed executing UpdateChannelStock query with buffer stock: %s", err.Error())
		}

		_, err = stmtRelatedChannelStock.Exec(item.Delta, input.WarehouseID, item.SKU, item.GateID, item.ChannelID)
		if err != nil {
			return fmt.Errorf("failed executing UpdateChannelStock query for related stock: %s", err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed committing UpdateChannelStock transaction: %s", err.Error())
	}

	return nil
}
