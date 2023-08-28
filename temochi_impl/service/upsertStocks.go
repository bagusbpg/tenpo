package service

import (
	"context"
	"fmt"

	"github.com/bagusbpg/tenpo/temochi"
)

func (ths *service) UpsertStocks(ctx context.Context, req temochi.UpsertStocksReq, res *temochi.UpsertStocksRes) error {
	input, failedSpecs := constructUpsertStockInput(req)
	if len(failedSpecs) > 0 {
		*res = temochi.UpsertStocksRes{
			FailedSpecs: failedSpecs,
		}
		return fmt.Errorf("failed validating channel stock")
	}

	err := ths.repository.UpsertStock(ctx, input, nil)
	if err != nil {
		*res = temochi.UpsertStocksRes{
			FailedSpecs: append(failedSpecs, constructFailedSpecsFromInputDB(input)...),
		}
		return fmt.Errorf("failed at repository.UpsertStock: %s", err.Error())
	}

	*res = temochi.UpsertStocksRes{
		FailedSpecs: failedSpecs,
	}

	return nil
}

func constructUpsertStockInput(req temochi.UpsertStocksReq) (UpsertStockDBInput, []temochi.FailedUpsertStockSpec) {
	input := UpsertStockDBInput{WarehouseID: req.WarehouseID}
	failedSpecs := make([]temochi.FailedUpsertStockSpec, 0, 10)
	for i := range req.UpsertStockSpecs {
		if valid, message := validateChannelStock(req.UpsertStockSpecs[i]); valid {
			inventoryInput := UpsertInventoryInput{
				SKU:         req.UpsertStockSpecs[i].SKU,
				Stock:       req.UpsertStockSpecs[i].Stock,
				BufferStock: req.UpsertStockSpecs[i].BufferStock,
			}
			upsertChannelStockInputs := make([]UpsertChannelStockInput, 0, 10)
			for j := range req.UpsertStockSpecs[i].ChannelStockSpecs {
				upsertChannelStockInputs = append(upsertChannelStockInputs, UpsertChannelStockInput{
					SKU:       req.UpsertStockSpecs[i].SKU,
					GateID:    req.UpsertStockSpecs[i].ChannelStockSpecs[j].GateID,
					ChannelID: req.UpsertStockSpecs[i].ChannelStockSpecs[j].ChannelID,
					Stock:     req.UpsertStockSpecs[i].ChannelStockSpecs[j].Stock,
				})
			}
			input.UpsertInventoryInputs = append(input.UpsertInventoryInputs, inventoryInput)
			input.UpsertChannelStockInputs = append(input.UpsertChannelStockInputs, upsertChannelStockInputs...)
		} else {
			failedSpecs = append(failedSpecs, temochi.FailedUpsertStockSpec{
				SKU:     req.UpsertStockSpecs[i].SKU,
				Message: message,
			})
		}
	}

	return input, failedSpecs
}

func validateChannelStock(spec temochi.UpsertStockSpec) (bool, string) {
	for i := range spec.ChannelStockSpecs {
		if spec.ChannelStockSpecs[i].Stock > spec.Stock-spec.BufferStock {
			return false, fmt.Sprintf("channel stock at gateID <%s> and channelID <%s> cannot be greater than available stock", spec.ChannelStockSpecs[i].GateID, spec.ChannelStockSpecs[i].ChannelID)
		}
	}
	return true, ""
}

func constructFailedSpecsFromInputDB(input UpsertStockDBInput) []temochi.FailedUpsertStockSpec {
	failedSpecs := make([]temochi.FailedUpsertStockSpec, 0, 10)
	for i := range input.UpsertInventoryInputs {
		failedSpecs = append(failedSpecs, temochi.FailedUpsertStockSpec{
			SKU:     input.UpsertInventoryInputs[i].SKU,
			Message: "failed upserting stock",
		})
	}

	return failedSpecs
}
