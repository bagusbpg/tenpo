package service

import (
	"context"
	"fmt"

	"github.com/bagusbpg/tenpo/temochi"
	"github.com/bagusbpg/tenpo/temochi_impl/repository"
)

func (ths *service) GetStocks(ctx context.Context, req temochi.GetStocksReq, res *temochi.GetStocksRes) error {
	if len(req.SKUs) == 0 {
		return nil
	}

	input := repository.GetStocksDBInput{
		WarehouseID: req.WarehouseID,
		SKUs:        req.SKUs,
	}
	output := repository.GetStocksDBOutput{}
	if err := ths.repository.GetStocks(ctx, input, &output); err != nil {
		return fmt.Errorf("failed at repository.GetStocks: %s", err.Error())
	}

	for _, stock := range output.Stocks {
		res.Stocks = append(res.Stocks, temochi.Stock{
			Inventory:     stock.Inventory,
			ChannelStocks: stock.ChannelStocks,
		})
	}

	return nil
}
