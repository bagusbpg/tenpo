package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bagusbpg/tenpo/temochi"
	"github.com/bagusbpg/tenpo/temochi/mock"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
)

func TestUpdateChannelStocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedService := mock.NewMockService(ctrl)
	testHandler := New(mockedService, validator.New())

	req := temochi.UpdateChannelStocksReq{
		ActorID:   "dummy-actor-id",
		ActorName: "dummy-actor-name",
		UpdateChannelStockSpecs: []temochi.UpdateChannelStockSpec{
			{
				SKU:       "dummy-sku",
				GateID:    "dummy-gate-id",
				ChannelID: "dummy-channel-id",
				Delta:     -1,
			},
		},
	}

	t.Run("With invalid JSON as request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/stocks/dummy-warehouse-id", bytes.NewReader([]byte("INVALID_JSON")))

		testHandler.UpdateChannelStocks().ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Error("response status code should be bad request")
		}
	})

	t.Run("With failure at struct validation", func(t *testing.T) {
		reqByte, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/stocks/dummy-warehouse-id", bytes.NewReader(reqByte))

		testHandler.UpdateChannelStocks().ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Error("response status code should be bad request")
		}
	})

	t.Run("With param's warehouse_id != request body's warehouse_id", func(t *testing.T) {
		req.WarehouseID = "invalid-warehouse-id"
		reqByte, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/stocks/dummy-warehouse-id", bytes.NewReader(reqByte))

		testHandler.UpdateChannelStocks().ServeHTTP(w, r)

		if w.Code != http.StatusForbidden {
			t.Error("response status code should be forbidden")
		}
	})

	req.WarehouseID = "dummy-warehouse-id"

	t.Run("With service.UpdateChannelStocks returns error", func(t *testing.T) {
		mockedService.EXPECT().
			UpdateChannelStocks(context.Background(), req, nil).
			Return(errors.New("dummy-error"))

		reqByte, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/stocks/dummy-warehouse-id", bytes.NewReader(reqByte))

		testHandler.UpdateChannelStocks().ServeHTTP(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Error("response status code should be internal server error")
		}
	})

	t.Run("With service.UpdateChannelStocks returns no error", func(t *testing.T) {
		mockedService.EXPECT().
			UpdateChannelStocks(context.Background(), req, nil).
			Return(nil)

		reqByte, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPatch, "/stocks/dummy-warehouse-id", bytes.NewReader(reqByte))

		testHandler.UpdateChannelStocks().ServeHTTP(w, r)

		if w.Code != http.StatusNoContent {
			t.Error("response status code should be no content")
		}
	})
}
