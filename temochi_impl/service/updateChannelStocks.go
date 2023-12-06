package service

import (
	"context"
	"fmt"

	"github.com/bagusbpg/tenpo/temochi"
	repository "github.com/bagusbpg/tenpo/temochi_impl/repository"
)

func (ths service) UpdateChannelStocks(ctx context.Context, req temochi.UpdateChannelStocksReq, res *temochi.UpdateChannelStocksRes) error {
	input := repository.UpdateChannelStocksDBInput{WarehouseID: req.WarehouseID}
	for i := range req.UpdateChannelStockSpecs {
		input.UpdateChannelStockInputs = append(input.UpdateChannelStockInputs, repository.UpdateChannelStockInput{
			SKU:       req.UpdateChannelStockSpecs[i].SKU,
			GateID:    req.UpdateChannelStockSpecs[i].GateID,
			ChannelID: req.UpdateChannelStockSpecs[i].ChannelID,
			Delta:     req.UpdateChannelStockSpecs[i].Delta,
		})
	}

	err := ths.repository.UpdateChannelStocks(ctx, input, nil)
	if err != nil {
		return fmt.Errorf("failed at repository.UpdateChannelStocks: %s", err.Error())
	}

	return nil
}
