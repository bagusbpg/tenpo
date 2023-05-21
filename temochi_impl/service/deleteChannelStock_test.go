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

func TestDeleteChannelStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedRepository := mock.NewMockRepository(ctrl)
	testService := New(mockedRepository)

	ctx := context.TODO()
	req := temochi.DeleteChannelStockReq{
		WarehouseID: "dummy-warehouse-id",
		GateID:      "dummy-gate-id",
		ChannelID:   "dummy-channel-id",
	}
	input := repository.DeleteChannelStockDBInput{
		WarehouseID: "dummy-warehouse-id",
		GateID:      "dummy-gate-id",
		ChannelID:   "dummy-channel-id",
	}

	t.Run("With repository.DeleteChannelStock returns error", func(t *testing.T) {
		mockedRepository.EXPECT().DeleteChannelStock(ctx, input, nil).Return(errors.New("dummy-error"))

		err := testService.DeleteChannelStock(ctx, req, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With repository.DeleteChannelStock returns no error", func(t *testing.T) {
		mockedRepository.EXPECT().DeleteChannelStock(ctx, input, nil).Return(nil)

		err := testService.DeleteChannelStock(ctx, req, nil)
		if err != nil {
			t.Error("nil should be returned")
		}
	})
}
