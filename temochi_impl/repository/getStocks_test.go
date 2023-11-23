package repository

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bagusbpg/tenpo/temochi"
)

func TestGetStocks(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed opening a stub database connection: %s", err.Error())
	}
	defer db.Close()

	mockedRepository := New(db)

	ctx := context.TODO()
	input := GetStocksDBInput{
		WarehouseID: "dummy-warehouse-id",
		SKUs:        []string{"dummy-sku-1", "dummy-sku-2"},
	}
	output := []StockDB{{Inventory: temochi.Inventory{WarehouseID: input.WarehouseID}}}
	outputByte, _ := json.Marshal(output)

	query, _ := buildGetStocksQuery(input)

	t.Run("With ExecContext returns error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(input.WarehouseID, input.SKUs[0], input.SKUs[1], input.WarehouseID, input.SKUs[0], input.SKUs[1]).
			WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.GetStocks(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With ExecContext returns no error", func(t *testing.T) {
		t.Run("But returned row is nil", func(t *testing.T) {
			nilRows := sqlmock.NewRows([]string{"result_json"})

			mock.
				ExpectQuery(query).
				WithArgs(input.WarehouseID, input.SKUs[0], input.SKUs[1], input.WarehouseID, input.SKUs[0], input.SKUs[1]).
				WillReturnRows(nilRows)

			err := mockedRepository.GetStocks(ctx, input, nil)
			if err != nil {
				t.Error("nil should be returned")
			}
		})

		t.Run("But returned row is invalid JSON", func(t *testing.T) {
			invalidJSONRows := sqlmock.NewRows([]string{"result_json"}).AddRow("INVALID_JSON")

			mock.
				ExpectQuery(query).
				WithArgs(input.WarehouseID, input.SKUs[0], input.SKUs[1], input.WarehouseID, input.SKUs[0], input.SKUs[1]).
				WillReturnRows(invalidJSONRows)

			output := new(GetStocksDBOutput)
			err := mockedRepository.GetStocks(ctx, input, output)
			if err == nil {
				t.Error("error should be returned")
			}
		})

		t.Run("And returned row is valid JSON", func(t *testing.T) {
			validJSONRows := sqlmock.NewRows([]string{"result_json"}).AddRow(string(outputByte))

			mock.
				ExpectQuery(query).
				WithArgs(input.WarehouseID, input.SKUs[0], input.SKUs[1], input.WarehouseID, input.SKUs[0], input.SKUs[1]).
				WillReturnRows(validJSONRows)

			output := new(GetStocksDBOutput)
			err := mockedRepository.GetStocks(ctx, input, output)
			if err != nil {
				t.Error("nil should be returned")
			}
		})
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
