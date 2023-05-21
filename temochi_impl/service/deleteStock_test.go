package service

import (
	"context"
	"errors"
	"testing"

	"github.com/bagusbpg/tenpo/temochi"
	"github.com/bagusbpg/tenpo/temochi_impl/repository"
	"github.com/bagusbpg/tenpo/temochi_impl/repository/mock"
	"github.com/golang/mock/gomock"
)

func TestDeleteStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedRepository := mock.NewMockRepository(ctrl)
	testService := New(mockedRepository)

	ctx := context.TODO()
	req := temochi.DeleteStockReq{
		WarehouseID: "dummy-warehouse-id",
		SKU:         "dummy-sku",
	}
	input := repository.DeleteStockDBInput{
		WarehouseID: "dummy-warehouse-id",
		SKU:         "dummy-sku",
	}

	t.Run("With repository.DeleteStock returns error", func(t *testing.T) {
		mockedRepository.EXPECT().DeleteStock(ctx, input, nil).Return(errors.New("dummy-error"))

		err := testService.DeleteStock(ctx, req, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With repository.DeleteStock returns no error", func(t *testing.T) {
		mockedRepository.EXPECT().DeleteStock(ctx, input, nil).Return(nil)

		err := testService.DeleteStock(ctx, req, nil)
		if err != nil {
			t.Error("nil should be returned")
		}
	})
}
