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

func TestUpsertStocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedService := mock.NewMockService(ctrl)
	testHandler := New(mockedService, validator.New())

	req := temochi.UpsertStocksReq{
		ActorID:   "dummy-actor-id",
		ActorName: "dummy-actor-name",
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

	t.Run("With invalid JSON as request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/stocks/dummy-warehouse-id", bytes.NewReader([]byte("INVALID_JSON")))

		testHandler.UpsertStocks().ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Error("response status code should be bad request")
		}
	})

	t.Run("With failure at struct validation", func(t *testing.T) {
		reqByte, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/stocks/dummy-warehouse-id", bytes.NewReader(reqByte))

		testHandler.UpsertStocks().ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Error("response status code should be bad request")
		}
	})

	t.Run("With param's warehouse_id != request body's warehouse_id", func(t *testing.T) {
		req.WarehouseID = "invalid-warehouse-id"
		reqByte, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/stocks/dummy-warehouse-id", bytes.NewReader(reqByte))

		testHandler.UpsertStocks().ServeHTTP(w, r)

		if w.Code != http.StatusForbidden {
			t.Error("response status code should be forbidden")
		}
	})

	req.WarehouseID = "dummy-warehouse-id"

	t.Run("With service.UpdateChannelStocks returns error on validating channel stock", func(t *testing.T) {
		res := temochi.UpsertStocksRes{}
		mockedService.EXPECT().
			UpsertStocks(context.Background(), req, &res).
			Return(errors.New("failed to validate channel stock"))

		reqByte, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/stocks/dummy-warehouse-id", bytes.NewReader(reqByte))

		testHandler.UpsertStocks().ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Error("response status code should be bad request")
		}
	})

	req.UpsertStockSpecs[0].ChannelStockSpecs[0].Stock = 7

	t.Run("With service.UpdateChannelStocks returns error from repository", func(t *testing.T) {
		res := temochi.UpsertStocksRes{}
		mockedService.EXPECT().
			UpsertStocks(context.Background(), req, &res).
			Return(errors.New("dummy-error"))

		reqByte, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/stocks/dummy-warehouse-id", bytes.NewReader(reqByte))

		testHandler.UpsertStocks().ServeHTTP(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Error("response status code should be internal server error")
		}
	})

	t.Run("failed at repository.UpserStock", func(t *testing.T) {
		res := temochi.UpsertStocksRes{}
		mockedService.EXPECT().
			UpsertStocks(context.Background(), req, &res).
			Return(nil)

		reqByte, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/stocks/dummy-warehouse-id", bytes.NewReader(reqByte))

		testHandler.UpsertStocks().ServeHTTP(w, r)

		if w.Code != http.StatusCreated {
			t.Error("response status code should be created")
		}
	})
}
