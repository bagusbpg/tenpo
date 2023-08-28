package repository

import (
	"context"
	"fmt"

	"github.com/bagusbpg/tenpo/temochi_impl/service"
)

const DELETE_STOCK_QUERY = `
DELETE FROM "temochi".inventory
WHERE warehouse_id = $1 AND sku = $2`

func (ths *repository) DeleteStock(ctx context.Context, input service.DeleteStockDBInput, output *service.DeleteStockDBOutput) error {
	_, err := ths.db.ExecContext(ctx, DELETE_STOCK_QUERY, input.WarehouseID, input.SKU)
	if err != nil {
		return fmt.Errorf("failed executing query: %s", err.Error())
	}

	return nil
}
