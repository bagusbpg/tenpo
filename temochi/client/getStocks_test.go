package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bagusbpg/tenpo/temochi"
)

func TestGetStocks(t *testing.T) {
	req := temochi.GetStocksReq{
		WarehouseID: "dummy-warehouse-id",
		SKUs:        []string{"dummy-sku"},
	}
	res := temochi.GetStocksRes{
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
	resByte, _ := json.Marshal(res)

	t.Run("When res is nil", func(t *testing.T) {
		client := NewTemochiClient(Config{})

		err := client.GetStocks(context.TODO(), req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("When failed parsing request URI", func(t *testing.T) {
		client := NewTemochiClient(Config{})

		err := client.GetStocks(context.TODO(), req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("When failed creating request", func(t *testing.T) {
		client := NewTemochiClient(Config{URL: "/0.0.0.0:8000"})

		var nilCtx context.Context
		err := client.GetStocks(nilCtx, req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("When failed sending request", func(t *testing.T) {
		client := NewTemochiClient(Config{URL: "/0.0.0.0:8000"})

		err := client.GetStocks(context.TODO(), req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("When response's status code is below 200 or greater than equal to 300", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer testServer.Close()

		client := NewTemochiClient(Config{URL: testServer.URL})

		err := client.GetStocks(context.TODO(), req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("When parsing response body returns error", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer testServer.Close()

		client := NewTemochiClient(Config{URL: testServer.URL})

		err := client.GetStocks(context.TODO(), req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("When parsing response body returns no error", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(resByte)
		}))
		defer testServer.Close()

		client := NewTemochiClient(Config{URL: testServer.URL})

		err := client.GetStocks(context.TODO(), req, &res)
		if err != nil {
			t.Error("nil should be returned")
		}
	})
}
