package service

import (
	"context"
	"fmt"

	"github.com/bagusbpg/tenpo/temochi"
)

func (ths *service) DeleteStock(ctx context.Context, req temochi.DeleteStockReq, res *temochi.DeleteStockRes) error {
	input := DeleteStockDBInput{
		WarehouseID: req.WarehouseID,
		SKU:         req.SKU,
	}
	err := ths.repository.DeleteStock(ctx, input, nil)
	if err != nil {
		return fmt.Errorf("failed at repository.DeleteStock: %s", err.Error())
	}

	return nil
}
