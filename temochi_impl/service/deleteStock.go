package service

import (
	"context"
	"fmt"

	"github.com/bagusbpg/tenpo/temochi"
	repository "github.com/bagusbpg/tenpo/temochi_impl/repository"
)

func (ths *service) DeleteStock(ctx context.Context, req temochi.DeleteStockReq, res *temochi.DeleteStockRes) error {
	input := repository.DeleteStockDBInput{
		WarehouseID: req.WarehouseID,
		SKU:         req.SKU,
	}
	err := ths.repository.DeleteStock(ctx, input, nil)
	if err != nil {
		return fmt.Errorf("failed at repository.DeleteStock: %s", err.Error())
	}

	return nil
}
