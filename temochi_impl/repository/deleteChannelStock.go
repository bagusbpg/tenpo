package repository

import (
	"context"
	"fmt"
)

type DeleteChannelStockDBInput struct {
	WarehouseID string
	GateID      string
	ChannelID   string
}

type DeleteChannelStockDBOutput struct{}

const DELETE_CHANNEL_STOCK_QUERY = `
DELETE FROM "temochi".channel_stock
WHERE warehouse_id = $1 AND gate_id = $2 AND channel_id = $3`

func (ths repository) DeleteChannelStock(ctx context.Context, input DeleteChannelStockDBInput, res *DeleteChannelStockDBOutput) error {
	_, err := ths.db.ExecContext(ctx, DELETE_CHANNEL_STOCK_QUERY, input.WarehouseID, input.GateID, input.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %s", err.Error())
	}

	return nil
}
