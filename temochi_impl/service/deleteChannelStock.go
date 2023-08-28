package service

import (
	"context"
	"fmt"

	"github.com/bagusbpg/tenpo/temochi"
)

func (ths *service) DeleteChannelStock(ctx context.Context, req temochi.DeleteChannelStockReq, res *temochi.DeleteChannelStockRes) error {
	input := DeleteChannelStockDBInput{
		WarehouseID: req.WarehouseID,
		GateID:      req.GateID,
		ChannelID:   req.ChannelID,
	}
	err := ths.repository.DeleteChannelStock(ctx, input, nil)
	if err != nil {
		return fmt.Errorf("failed at repository.DeleteChannelStock: %s", err.Error())
	}

	return nil
}
