package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bagusbpg/tenpo/temochi"
	"github.com/bagusbpg/tenpo/temochi/mock"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
)

func TestDeleteStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedService := mock.NewMockService(ctrl)
	testHandler := New(mockedService, validator.New())

	req := temochi.DeleteStockReq{
		WarehouseID: "dummy-warehouse-id",
		SKU:         "dummy-sku",
	}

	t.Run("With service.DeleteStock returns error", func(t *testing.T) {
		mockedService.EXPECT().
			DeleteStock(context.Background(), req, nil).
			Return(errors.New("dummy-error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/stocks/"+req.WarehouseID+"/"+req.SKU, nil)

		testHandler.DeleteStock().ServeHTTP(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Error("response status code should be internal server error", w.Code)
		}
	})

	t.Run("With service.DeleteStock returns no error", func(t *testing.T) {
		mockedService.EXPECT().
			DeleteStock(context.Background(), req, nil).
			Return(nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/stocks/"+req.WarehouseID+"/"+req.SKU, nil)

		testHandler.DeleteStock().ServeHTTP(w, r)

		if w.Code != http.StatusNoContent {
			t.Error("response status code should be no content", w.Code)
		}
	})
}
