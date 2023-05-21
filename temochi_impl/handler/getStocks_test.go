package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bagusbpg/tenpo/temochi"
	"github.com/bagusbpg/tenpo/temochi/mock"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
)

func TestGetStocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedService := mock.NewMockService(ctrl)
	testHandler := New(mockedService, validator.New())

	req := temochi.GetStocksReq{
		WarehouseID: "dummy-warehouse-id",
		SKUs:        []string{"dummy-sku"},
	}

	t.Run("With ParseQuery returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/stocks/dummy-warehouse-id?skus=dum%3Bmy-sku", nil)

		testHandler.GetStocks().ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Error("response status code should be bad request")
		}
	})

	t.Run("With service.GetStocks returns error", func(t *testing.T) {
		mockedService.EXPECT().
			GetStocks(context.Background(), req, new(temochi.GetStocksRes)).
			Return(errors.New("dummy-error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/stocks/dummy-warehouse-id?skus=dummy-sku", nil)

		testHandler.GetStocks().ServeHTTP(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Error("response status code should be internal server error")
		}
	})

	t.Run("With service.GetStocks returns error", func(t *testing.T) {
		res := temochi.GetStocksRes{}
		mockedService.EXPECT().
			GetStocks(context.Background(), req, &res).
			Do(func(_ context.Context, _ temochi.GetStocksReq, res *temochi.GetStocksRes) {
				*res = temochi.GetStocksRes{
					Stocks: []temochi.Stock{
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
					},
				}
			}).
			Return(nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/stocks/dummy-warehouse-id?skus=dummy-sku", nil)

		testHandler.GetStocks().ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Error("response status code should be ok")
		}
	})
}
