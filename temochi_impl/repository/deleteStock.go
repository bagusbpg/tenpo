package repository

import (
	"context"
	"fmt"
)

type DeleteStockDBInput struct {
	WarehouseID string
	SKU         string
}

type DeleteStockDBOutput struct{}

const DELETE_STOCK_QUERY = `
DELETE FROM "temochi".inventory
WHERE warehouse_id = $1 AND sku = $2`

func (ths repository) DeleteStock(ctx context.Context, input DeleteStockDBInput, output *DeleteStockDBOutput) error {
	_, err := ths.db.ExecContext(ctx, DELETE_STOCK_QUERY, input.WarehouseID, input.SKU)
	if err != nil {
		return fmt.Errorf("failed to execute query: %s", err.Error())
	}

	return nil
}
