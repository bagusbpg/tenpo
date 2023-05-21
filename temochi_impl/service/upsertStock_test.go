package service

import (
	"context"
	"errors"
	"testing"

	"github.com/bagusbpg/tenpo/temochi"
	"github.com/bagusbpg/tenpo/temochi_impl/repository/mock"
	"github.com/golang/mock/gomock"
)

func TestUpsertStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedRepository := mock.NewMockRepository(ctrl)
	testService := New(mockedRepository)

	ctx := context.TODO()
	req := temochi.UpsertStocksReq{
		ActorID:     "dummy-actor-id",
		ActorName:   "dummy-actor-name",
		WarehouseID: "dummy-warehouse-id",
		UpsertStockSpecs: []temochi.UpsertStockSpec{
			{
				SKU:         "dummy-sku",
				Stock:       10,
				BufferStock: 2,
				ChannelStockSpecs: []temochi.ChannelStockSpec{
					{
						GateID:    "dummy-gate-id",
						ChannelID: "dummy-channel-id",
						Stock:     9,
					},
				},
			},
		},
	}

	t.Run("With all request are invalid", func(t *testing.T) {
		res := temochi.UpsertStocksRes{}
		err := testService.UpsertStocks(ctx, req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
		if len(res.FailedSpecs) == 0 {
			t.Error("invalid specs should be appended to failed specs")
		}
	})

	req.UpsertStockSpecs[0].ChannelStockSpecs[0].Stock = 7
	input, _ := constructUpsertStockInput(req)

	t.Run("With repository.UpsertStock returns error", func(t *testing.T) {
		mockedRepository.EXPECT().
			UpsertStock(ctx, input, nil).
			Return(errors.New("dummy-error"))

		res := temochi.UpsertStocksRes{}
		err := testService.UpsertStocks(ctx, req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
		if len(res.FailedSpecs) == 0 {
			t.Error("Unupserted specs should be appended to failed specs")
		}
	})

	t.Run("With repository.UpsertStock returns no error", func(t *testing.T) {
		mockedRepository.EXPECT().
			UpsertStock(ctx, input, nil).
			Return(nil)

		res := temochi.UpsertStocksRes{}
		err := testService.UpsertStocks(ctx, req, &res)
		if err != nil {
			t.Error("nil should be returned")
		}
		if len(res.FailedSpecs) != 0 {
			t.Error("failed specs must be empty")
		}
	})
}
