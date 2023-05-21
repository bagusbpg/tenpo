package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUpsertStocks(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed opening a stub database connection: %s", err.Error())
	}
	defer db.Close()

	mockedRepository := New(db)

	ctx := context.TODO()
	input := UpsertStockDBInput{
		WarehouseID:              "dummy-warehouse-id",
		UpsertInventoryInputs:    []UpsertInventoryInput{{SKU: "dummy-sku", Stock: 10, BufferStock: 2}},
		UpsertChannelStockInputs: []UpsertChannelStockInput{{SKU: "dummy-sku", GateID: "dummy-gate-id", ChannelID: "dummy-channel-id", Stock: 5}},
	}

	queryUpsertInventory, _ := buildUpsertInventoryQuery(input)
	queryUpsertChannelStock, _ := buildUpsertChannelStockQuery(input)
	queryDeleteChannelStock, _ := buildDeleteExcludedGateChannelStockQuery(input)

	t.Run("With Begin returns error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.UpsertStock(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With ExecContext queryUpsertInventory returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(queryUpsertInventory).
			WithArgs(input.WarehouseID, input.UpsertInventoryInputs[0].SKU, input.UpsertInventoryInputs[0].Stock, input.UpsertInventoryInputs[0].BufferStock).
			WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.UpsertStock(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With ExecContext queryUpsertChannelStock returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(queryUpsertInventory).
			WithArgs(input.WarehouseID, input.UpsertInventoryInputs[0].SKU, input.UpsertInventoryInputs[0].Stock, input.UpsertInventoryInputs[0].BufferStock).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(queryUpsertChannelStock).
			WithArgs(input.WarehouseID, input.UpsertChannelStockInputs[0].SKU, input.UpsertChannelStockInputs[0].GateID, input.UpsertChannelStockInputs[0].ChannelID, input.UpsertChannelStockInputs[0].Stock).
			WillReturnError(errors.New("dummy-error"))
		mock.ExpectRollback()

		err := mockedRepository.UpsertStock(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With ExecContext queryDeleteChannelStock returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(queryUpsertInventory).
			WithArgs(input.WarehouseID, input.UpsertInventoryInputs[0].SKU, input.UpsertInventoryInputs[0].Stock, input.UpsertInventoryInputs[0].BufferStock).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(queryUpsertChannelStock).
			WithArgs(input.WarehouseID, input.UpsertChannelStockInputs[0].SKU, input.UpsertChannelStockInputs[0].GateID, input.UpsertChannelStockInputs[0].ChannelID, input.UpsertChannelStockInputs[0].Stock).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(queryDeleteChannelStock).
			WithArgs(input.WarehouseID, input.UpsertChannelStockInputs[0].SKU, input.UpsertChannelStockInputs[0].SKU+"#"+input.UpsertChannelStockInputs[0].GateID+"#"+input.UpsertChannelStockInputs[0].ChannelID).
			WillReturnError(errors.New("dummy-error"))
		mock.ExpectRollback()

		err := mockedRepository.UpsertStock(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With Commit returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(queryUpsertInventory).
			WithArgs(input.WarehouseID, input.UpsertInventoryInputs[0].SKU, input.UpsertInventoryInputs[0].Stock, input.UpsertInventoryInputs[0].BufferStock).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(queryUpsertChannelStock).
			WithArgs(input.WarehouseID, input.UpsertChannelStockInputs[0].SKU, input.UpsertChannelStockInputs[0].GateID, input.UpsertChannelStockInputs[0].ChannelID, input.UpsertChannelStockInputs[0].Stock).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(queryDeleteChannelStock).
			WithArgs(input.WarehouseID, input.UpsertChannelStockInputs[0].SKU, input.UpsertChannelStockInputs[0].SKU+"#"+input.UpsertChannelStockInputs[0].GateID+"#"+input.UpsertChannelStockInputs[0].ChannelID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.UpsertStock(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With Commit returns no error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(queryUpsertInventory).
			WithArgs(input.WarehouseID, input.UpsertInventoryInputs[0].SKU, input.UpsertInventoryInputs[0].Stock, input.UpsertInventoryInputs[0].BufferStock).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(queryUpsertChannelStock).
			WithArgs(input.WarehouseID, input.UpsertChannelStockInputs[0].SKU, input.UpsertChannelStockInputs[0].GateID, input.UpsertChannelStockInputs[0].ChannelID, input.UpsertChannelStockInputs[0].Stock).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(queryDeleteChannelStock).
			WithArgs(input.WarehouseID, input.UpsertChannelStockInputs[0].SKU, input.UpsertChannelStockInputs[0].SKU+"#"+input.UpsertChannelStockInputs[0].GateID+"#"+input.UpsertChannelStockInputs[0].ChannelID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := mockedRepository.UpsertStock(ctx, input, nil)
		if err != nil {
			t.Error("nil should be returned")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
