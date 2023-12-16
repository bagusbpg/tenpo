package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDeleteStock(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %s", err.Error())
	}
	defer db.Close()

	mockedRepository := New(db)

	ctx := context.TODO()
	input := DeleteStockDBInput{
		WarehouseID: "dummy-warehouse-id",
		SKU:         "dummy-sku",
	}

	t.Run("With ExecContext returns error", func(t *testing.T) {
		mock.
			ExpectExec(DELETE_STOCK_QUERY).
			WithArgs(input.WarehouseID, input.SKU).
			WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.DeleteStock(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With ExecContext returns no error", func(t *testing.T) {
		mock.
			ExpectExec(DELETE_STOCK_QUERY).
			WithArgs(input.WarehouseID, input.SKU).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := mockedRepository.DeleteStock(ctx, input, nil)
		if err != nil {
			t.Error("nil should be returned")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
