package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagusbpg/tenpo/temochi_impl/service"
)

func (ths *repository) GetStocks(ctx context.Context, input service.GetStocksDBInput, output *service.GetStocksDBOutput) error {
	query, args := buildGetStocksQuery(ctx, input)

	var res sql.NullString
	err := ths.db.QueryRowContext(ctx, query, args...).Scan(&res)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed executing query: %s", err.Error())
	}

	if !res.Valid {
		return nil
	}

	err = json.Unmarshal([]byte(res.String), &output.Stocks)
	if err != nil {
		return fmt.Errorf("failed parsing result: %s", err.Error())
	}

	return nil
}

func buildGetStocksQuery(ctx context.Context, input service.GetStocksDBInput) (string, []interface{}) {
	channelStockCTEQuery, channelStockArgs := sq.
		Select(
			"warehouse_id",
			"sku",
			`ARRAY_AGG(JSON_BUILD_OBJECT('warehouseId', warehouse_id, 'sku', sku, 'gateId', gate_id, 'channelId', channel_id, 'stock', stock, 'version', version, 'createdAt', CAST(EXTRACT(EPOCH FROM created_at) AS integer), 'updatedAt', CAST(EXTRACT(EPOCH FROM updated_at) AS integer))) AS datarow`,
		).
		From(`"temochi".channel_stock`).
		Where(sq.And{
			sq.Eq{"warehouse_id": input.WarehouseID},
			sq.Eq{"sku": input.SKUs},
		}).
		GroupBy(
			"warehouse_id",
			"sku",
		).
		MustSql()

	queryResultCTE, queryResultArgs := sq.
		Select(
			"inventory.warehouse_id AS warehouseId",
			"inventory.sku",
			"inventory.stock",
			"inventory.buffer_stock AS bufferStock",
			"inventory.version",
			"CAST(EXTRACT(EPOCH FROM inventory.updated_at) AS integer) AS updatedAt",
			"CAST(EXTRACT(EPOCH FROM inventory.created_at) AS integer) AS createdAt",
			"channel_stock.datarow AS channelStocks",
		).
		From(`"temochi".inventory`).
		LeftJoin("channel_stock ON channel_stock.warehouse_id = inventory.warehouse_id AND channel_stock.sku = inventory.sku").
		Where(sq.And{
			sq.Eq{"inventory.warehouse_id": input.WarehouseID},
			sq.Eq{"inventory.sku": input.SKUs},
		}).
		MustSql()

	finalQueryBuilder := sq.
		Select("ARRAY_TO_JSON(ARRAY_AGG(TO_JSON(query_result))) AS result_json").
		Prefix(fmt.Sprintf("WITH channel_stock AS (%s), query_result AS (%s)", channelStockCTEQuery, queryResultCTE), append(channelStockArgs, queryResultArgs...)...).
		From("query_result")

	return finalQueryBuilder.PlaceholderFormat(sq.Dollar).MustSql()
}
