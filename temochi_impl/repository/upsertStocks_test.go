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
		WarehouseID: "dummy-warehouse-id",
		UpsertInventoryInputs: []UpsertInventoryInput{
			{SKU: "dummy-sku-1", Stock: 10, BufferStock: 2},
			{SKU: "dummy-sku-2", Stock: 9, BufferStock: 1},
		},
		UpsertChannelStockInputs: []UpsertChannelStockInput{
			{SKU: "dummy-sku-1", GateID: "dummy-gate-id", ChannelID: "dummy-channel-id", Stock: 5},
			{SKU: "dummy-sku-2", GateID: "dummy-gate-id", ChannelID: "dummy-channel-id", Stock: 5},
		},
	}

	query, _ := buildUpsertStocksQuery(input)

	t.Run("With ExecContext returns error", func(t *testing.T) {
		mock.
			ExpectExec(query).
			WithArgs(
				input.WarehouseID,
				input.UpsertChannelStockInputs[0].SKU,
				input.UpsertChannelStockInputs[0].GateID,
				input.UpsertChannelStockInputs[0].ChannelID,
				input.UpsertChannelStockInputs[0].Stock,
				input.WarehouseID,
				input.UpsertChannelStockInputs[1].SKU,
				input.UpsertChannelStockInputs[1].GateID,
				input.UpsertChannelStockInputs[1].ChannelID,
				input.UpsertChannelStockInputs[1].Stock,
				input.WarehouseID,
				input.UpsertInventoryInputs[0].SKU,
				input.UpsertInventoryInputs[1].SKU,
				concatenateSKUGateIDChannelID(input.UpsertChannelStockInputs[0]),
				concatenateSKUGateIDChannelID(input.UpsertChannelStockInputs[1]),
				input.WarehouseID,
				input.UpsertInventoryInputs[0].SKU,
				input.UpsertInventoryInputs[0].Stock,
				input.UpsertInventoryInputs[0].BufferStock,
				input.WarehouseID,
				input.UpsertInventoryInputs[1].SKU,
				input.UpsertInventoryInputs[1].Stock,
				input.UpsertInventoryInputs[1].BufferStock,
			).
			WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.UpsertStock(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With ExecContext returns no error", func(t *testing.T) {
		mock.
			ExpectExec(query).
			WithArgs(
				input.WarehouseID,
				input.UpsertChannelStockInputs[0].SKU,
				input.UpsertChannelStockInputs[0].GateID,
				input.UpsertChannelStockInputs[0].ChannelID,
				input.UpsertChannelStockInputs[0].Stock,
				input.WarehouseID,
				input.UpsertChannelStockInputs[1].SKU,
				input.UpsertChannelStockInputs[1].GateID,
				input.UpsertChannelStockInputs[1].ChannelID,
				input.UpsertChannelStockInputs[1].Stock,
				input.WarehouseID,
				input.UpsertInventoryInputs[0].SKU,
				input.UpsertInventoryInputs[1].SKU,
				concatenateSKUGateIDChannelID(input.UpsertChannelStockInputs[0]),
				concatenateSKUGateIDChannelID(input.UpsertChannelStockInputs[1]),
				input.WarehouseID,
				input.UpsertInventoryInputs[0].SKU,
				input.UpsertInventoryInputs[0].Stock,
				input.UpsertInventoryInputs[0].BufferStock,
				input.WarehouseID,
				input.UpsertInventoryInputs[1].SKU,
				input.UpsertInventoryInputs[1].Stock,
				input.UpsertInventoryInputs[1].BufferStock,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := mockedRepository.UpsertStock(ctx, input, nil)
		if err != nil {
			t.Error("nil should be returned")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
