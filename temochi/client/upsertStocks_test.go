package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bagusbpg/tenpo/temochi"
)

func TestUpsertStocks(t *testing.T) {
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
						Stock:     5,
					},
				},
			},
		},
	}
	res := temochi.UpsertStocksRes{}
	resByte, _ := json.Marshal(res)

	t.Run("When res is nil", func(t *testing.T) {
		client := NewTemochiClient(Config{})

		err := client.UpsertStocks(context.TODO(), req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("When failed creating request", func(t *testing.T) {
		client := NewTemochiClient(Config{URL: "/0.0.0.0:8000"})

		var nilCtx context.Context
		err := client.UpsertStocks(nilCtx, req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("When failed sending request", func(t *testing.T) {
		client := NewTemochiClient(Config{URL: "/0.0.0.0:8000"})

		err := client.UpsertStocks(context.TODO(), req, &res)
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

		err := client.UpsertStocks(context.TODO(), req, &res)
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

		err := client.UpsertStocks(context.TODO(), req, &res)
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

		err := client.UpsertStocks(context.TODO(), req, &res)
		if err != nil {
			t.Error("nil should be returned")
		}
	})
}
