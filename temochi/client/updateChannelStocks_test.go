package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bagusbpg/tenpo/temochi"
)

func TestUpdateChannelStocks(t *testing.T) {
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
	res := temochi.UpdateChannelStocksRes{}

	t.Run("When res is nil", func(t *testing.T) {
		client := NewTemochiClient(Config{})

		err := client.UpdateChannelStocks(context.TODO(), req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("When failed creating request", func(t *testing.T) {
		client := NewTemochiClient(Config{URL: "/0.0.0.0:8000"})

		var nilCtx context.Context
		err := client.UpdateChannelStocks(nilCtx, req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("When failed sending request", func(t *testing.T) {
		client := NewTemochiClient(Config{URL: "/0.0.0.0:8000"})

		err := client.UpdateChannelStocks(context.TODO(), req, &res)
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

		err := client.UpdateChannelStocks(context.TODO(), req, &res)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("When response's status code is 200s", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer testServer.Close()

		client := NewTemochiClient(Config{URL: testServer.URL})

		err := client.UpdateChannelStocks(context.TODO(), req, &res)
		if err != nil {
			t.Error("nil should be returned")
		}
	})
}
