package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUpdateChannelStocks(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %s", err.Error())
	}
	defer db.Close()

	mockedRepository := New(db)

	ctx := context.TODO()
	input := UpdateChannelStocksDBInput{
		WarehouseID:              "dummy-warehouse-id",
		UpdateChannelStockInputs: []UpdateChannelStockInput{{SKU: "dummy-sku", GateID: "dummy-gate-id", ChannelID: "dummy-channel-id", Delta: -1}},
	}

	t.Run("With Begin returns error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.UpdateChannelStocks(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With PrepareContext stmtChannelStock returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectPrepare(UPDATE_CHANNEL_STOCK_QUERY).
			WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.UpdateChannelStocks(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With PrepareContext stmtInventoryWithBufferStock returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(UPDATE_CHANNEL_STOCK_QUERY)
		mock.ExpectPrepare(UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK).
			WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.UpdateChannelStocks(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With PrepareContext stmtRelatedChannelStock returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(UPDATE_CHANNEL_STOCK_QUERY)
		mock.ExpectPrepare(UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK)
		mock.ExpectPrepare(UPDATE_RELATED_CHANNEL_STOCK_QUERY).
			WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.UpdateChannelStocks(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With Exec UpdateChannelStock returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(UPDATE_CHANNEL_STOCK_QUERY)
		mock.ExpectPrepare(UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK)
		mock.ExpectPrepare(UPDATE_RELATED_CHANNEL_STOCK_QUERY)
		mock.
			ExpectExec(UPDATE_CHANNEL_STOCK_QUERY).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU, input.UpdateChannelStockInputs[0].GateID, input.UpdateChannelStockInputs[0].ChannelID).
			WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.UpdateChannelStocks(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With Exec UpdateInventory with buffer stock returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(UPDATE_CHANNEL_STOCK_QUERY)
		mock.ExpectPrepare(UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK)
		mock.ExpectPrepare(UPDATE_RELATED_CHANNEL_STOCK_QUERY)
		mock.
			ExpectExec(UPDATE_CHANNEL_STOCK_QUERY).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU, input.UpdateChannelStockInputs[0].GateID, input.UpdateChannelStockInputs[0].ChannelID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU).
			WillReturnError(errors.New("dummy-error"))
		mock.ExpectRollback()

		err := mockedRepository.UpdateChannelStocks(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With Exec UpdateRelatedChannelStock returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(UPDATE_CHANNEL_STOCK_QUERY)
		mock.ExpectPrepare(UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK)
		mock.ExpectPrepare(UPDATE_RELATED_CHANNEL_STOCK_QUERY)
		mock.
			ExpectExec(UPDATE_CHANNEL_STOCK_QUERY).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU, input.UpdateChannelStockInputs[0].GateID, input.UpdateChannelStockInputs[0].ChannelID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(UPDATE_RELATED_CHANNEL_STOCK_QUERY).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU, input.UpdateChannelStockInputs[0].GateID, input.UpdateChannelStockInputs[0].ChannelID).
			WillReturnError(errors.New("dummy-error"))
		mock.ExpectRollback()

		err := mockedRepository.UpdateChannelStocks(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With Commit returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(UPDATE_CHANNEL_STOCK_QUERY)
		mock.ExpectPrepare(UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK)
		mock.ExpectPrepare(UPDATE_RELATED_CHANNEL_STOCK_QUERY)
		mock.
			ExpectExec(UPDATE_CHANNEL_STOCK_QUERY).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU, input.UpdateChannelStockInputs[0].GateID, input.UpdateChannelStockInputs[0].ChannelID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(UPDATE_RELATED_CHANNEL_STOCK_QUERY).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU, input.UpdateChannelStockInputs[0].GateID, input.UpdateChannelStockInputs[0].ChannelID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(errors.New("dummy-error"))

		err := mockedRepository.UpdateChannelStocks(ctx, input, nil)
		if err == nil {
			t.Error("error should be returned")
		}
	})

	t.Run("With Commit returns no error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare(UPDATE_CHANNEL_STOCK_QUERY)
		mock.ExpectPrepare(UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK)
		mock.ExpectPrepare(UPDATE_RELATED_CHANNEL_STOCK_QUERY)
		mock.
			ExpectExec(UPDATE_CHANNEL_STOCK_QUERY).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU, input.UpdateChannelStockInputs[0].GateID, input.UpdateChannelStockInputs[0].ChannelID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(UPDATE_INVENTORY_QUERY_WITH_BUFFER_STOCK).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.
			ExpectExec(UPDATE_RELATED_CHANNEL_STOCK_QUERY).
			WithArgs(input.UpdateChannelStockInputs[0].Delta, input.WarehouseID, input.UpdateChannelStockInputs[0].SKU, input.UpdateChannelStockInputs[0].GateID, input.UpdateChannelStockInputs[0].ChannelID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := mockedRepository.UpdateChannelStocks(ctx, input, nil)
		if err != nil {
			t.Error("nil should be returned")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
