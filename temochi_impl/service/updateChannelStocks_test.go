package service

import (
	"context"
	"errors"
	"testing"

	"github.com/bagusbpg/tenpo/temochi"
	repository "github.com/bagusbpg/tenpo/temochi_impl/repository"
	"github.com/golang/mock/gomock"
)

func TestUpdateChannelStocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedRepository := NewMockRepository(ctrl)
	testService := New(mockedRepository)

	ctx := context.TODO()
	req := temochi.UpdateChannelStocksReq{
		ActorID:     "dummy-actor-id",
		ActorName:   "dummy-actor-name",
		WarehouseID: "dummy-warehouse-id",
		UpdateChannelStockSpecs: []temochi.UpdateChannelStockSpec{
			{
				SKU:       "dummy-sku",
				GateID:    "dummy-gate-id",
				ChannelID: "dummy-channel-id",
				Delta:     -1,
			},
		},
	}
	input := repository.UpdateChannelStocksDBInput{
		WarehouseID: req.WarehouseID,
		UpdateChannelStockInputs: []repository.UpdateChannelStockInput{
			{
				SKU:       req.UpdateChannelStockSpecs[0].SKU,
				GateID:    req.UpdateChannelStockSpecs[0].GateID,
				ChannelID: req.UpdateChannelStockSpecs[0].ChannelID,
				Delta:     req.UpdateChannelStockSpecs[0].Delta,
			},
		},
	}

	t.Run("With repository.UpdateChannelStocks returns error", func(t *testing.T) {
		mockedRepository.EXPECT().
			UpdateChannelStocks(ctx, input, nil).
			Return(errors.New("dummy-error"))

		err := testService.UpdateChannelStocks(ctx, req, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With repository.UpdateChannelStocks returns no error", func(t *testing.T) {
		mockedRepository.EXPECT().
			UpdateChannelStocks(ctx, input, nil).
			Return(nil)

		err := testService.UpdateChannelStocks(ctx, req, nil)
		if err != nil {
			t.Error("nil should be returned")
		}
	})
}
