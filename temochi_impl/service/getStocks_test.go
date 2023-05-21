package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bagusbpg/tenpo/temochi"
	"github.com/bagusbpg/tenpo/temochi_impl/repository"
	"github.com/bagusbpg/tenpo/temochi_impl/repository/mock"
	"github.com/golang/mock/gomock"
)

func TestGetStocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedRepository := mock.NewMockRepository(ctrl)
	testService := New(mockedRepository)

	ctx := context.TODO()
	req := temochi.GetStocksReq{}

	t.Run("With SKUs in request is empty", func(t *testing.T) {
		err := testService.GetStocks(ctx, req, nil)

		if err != nil {
			t.Error("nil should be returned")
		}
	})

	req.WarehouseID = "dummy-warehouse-id"
	req.SKUs = []string{"dummy-sku-1", "dummy-sku-2"}
	input := repository.GetStocksDBInput{
		WarehouseID: req.WarehouseID,
		SKUs:        req.SKUs,
	}
	output := repository.GetStocksDBOutput{}

	t.Run("With repository.GetStocks returns error", func(t *testing.T) {
		mockedRepository.EXPECT().
			GetStocks(ctx, input, &output).
			Return(errors.New("dummy-error"))

		err := testService.GetStocks(ctx, req, nil)

		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With repository.GetStocks returns no error", func(t *testing.T) {
		mockedRepository.EXPECT().
			GetStocks(ctx, input, &output).
			Do(func(_ context.Context, _ repository.GetStocksDBInput, output *repository.GetStocksDBOutput) {
				output.Stocks = []repository.StockDB{
					{
						Inventory: temochi.Inventory{
							WarehouseID: "dummy-warehouse-id",
							SKU:         "dummy-sku",
							Stock:       10,
							BufferStock: 2,
							Version:     0,
							CreatedAt:   time.Date(2015, 5, 16, 10, 0, 0, 0, time.Local).Unix(),
							UpdatedAt:   time.Date(2015, 5, 16, 10, 0, 0, 0, time.Local).Unix(),
						},
						ChannelStocks: []temochi.ChannelStock{
							{
								WarehouseID: "dummy-warehouse-id",
								SKU:         "dummy-sku",
								GateID:      "dummy-gate-id",
								ChannelID:   "dummy-channel-id",
								Stock:       5,
								Version:     0,
								CreatedAt:   time.Date(2015, 5, 16, 10, 0, 0, 0, time.Local).Unix(),
								UpdatedAt:   time.Date(2015, 5, 16, 10, 0, 0, 0, time.Local).Unix(),
							},
						},
					},
				}
			}).
			Return(nil)

		res := temochi.GetStocksRes{}
		err := testService.GetStocks(ctx, req, &res)

		if err != nil {
			t.Error("nil should be returned")
		}

		if len(res.Stocks) == 0 {
			t.Error("repository output must be appended to result")
		}
	})
}
