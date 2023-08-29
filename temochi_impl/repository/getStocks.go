package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/bagusbpg/tenpo/temochi_impl/service"
)

func (ths *repository) GetStocks(ctx context.Context, input service.GetStocksDBInput, output *service.GetStocksDBOutput) error {
	query, args := buildGetStocksQuery(input)

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

func buildGetStocksQuery(input service.GetStocksDBInput) (string, []interface{}) {
	query := strings.Builder{}
	args := make([]interface{}, 0)

	query.WriteString(`
		WITH channel_stock AS (
			SELECT
				warehouse_id,
				sku,
				ARRAY_AGG (
					JSON_BUILD_OBJECT (
						'warehouseId', warehouse_id,
						'sku', sku,
						'gateId', gate_id,
						'channelId', channel_id,
						'stock', stock,
						'version', version,
						'createdAt', CAST ( EXTRACT ( EPOCH FROM created_at ) AS integer ),
						'updatedAt', CAST ( EXTRACT ( EPOCH FROM updated_at ) AS integer )
					)
				) AS datarow
			FROM
				"temochi".channel_stock
			WHERE
	`)

	args = append(args, input.WarehouseID)
	query.WriteString(`warehouse_id = $` + strconv.Itoa(len(args)) + ` `)

	query.WriteString(`AND sku IN ( `)
	for i := range input.SKUs {
		args = append(args, input.SKUs[i])
		query.WriteString(`$` + strconv.Itoa(len(args)))
		if i < len(input.SKUs)-1 {
			query.WriteString(`, `)
		} else {
			query.WriteString(`) `)
		}
	}

	query.WriteString(`
			GROUP BY
				warehouse_id, sku
		),
		query_result AS (
			SELECT
				inventory.warehouse_id AS warehouseId,
				inventory.sku,
				inventory.stock,
				inventory.buffer_stock AS bufferStock,
				inventory.version,
				CAST ( EXTRACT ( EPOCH FROM inventory.updated_at ) AS integer ) AS updatedAt,
				CAST ( EXTRACT ( EPOCH FROM inventory.created_at ) AS integer ) AS createdAt,
				channel_stock.datarow AS channelStocks
			FROM
				"temochi".inventory
				LEFT JOIN channel_stock ON
					channel_stock.warehouse_id = inventory.warehouse_id
					AND channel_stock.sku = inventory.sku
			WHERE
	`)

	args = append(args, input.WarehouseID)
	query.WriteString(`inventory.warehouse_id = $` + strconv.Itoa(len(args)) + ` `)

	query.WriteString(`AND inventory.sku IN ( `)
	for i := range input.SKUs {
		args = append(args, input.SKUs[i])
		query.WriteString(`$` + strconv.Itoa(len(args)))
		if i < len(input.SKUs)-1 {
			query.WriteString(`, `)
		} else {
			query.WriteString(`) `)
		}
	}

	query.WriteString(`
		)
		SELECT
			ARRAY_TO_JSON (
				ARRAY_AGG (
					TO_JSON ( query_result )
				)
			) AS result_json
		FROM
			query_result
	`)

	return whitespaceNormalizer.ReplaceAllString(query.String(), " "), args
}
