package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bagusbpg/tenpo/temochi_impl/service"
)

func (ths *repository) UpsertStock(ctx context.Context, input service.UpsertStockDBInput, output *service.UpsertStockDBOutput) error {
	query, args := buildUpsertStocksQuery(input)
	fmt.Printf("query: %v\n", query)
	fmt.Printf("args: %v\n", args)

	_, err := ths.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed executing UpsertInventory query: %s", err.Error())
	}

	return nil
}

func buildUpsertStocksQuery(input service.UpsertStockDBInput) (string, []interface{}) {
	query := strings.Builder{}
	args := make([]interface{}, 0)

	if len(input.UpsertChannelStockInputs) > 0 {
		query.WriteString(`
			WITH upsert_channel_stock AS (
				INSERT INTO "temochi".channel_stock
					(
						warehouse_id,
						sku,
						gate_id,
						channel_id,
						stock
					)
				VALUES
		`)

		for i := range input.UpsertChannelStockInputs {
			query.WriteString(`( `)

			args = append(args, input.WarehouseID)
			query.WriteString(`$` + strconv.Itoa(len(args)) + `, `)

			args = append(args, input.UpsertChannelStockInputs[i].SKU)
			query.WriteString(`$` + strconv.Itoa(len(args)) + `, `)

			args = append(args, input.UpsertChannelStockInputs[i].GateID)
			query.WriteString(`$` + strconv.Itoa(len(args)) + `, `)

			args = append(args, input.UpsertChannelStockInputs[i].ChannelID)
			query.WriteString(`$` + strconv.Itoa(len(args)) + `, `)

			args = append(args, input.UpsertChannelStockInputs[i].Stock)
			query.WriteString(`$` + strconv.Itoa(len(args)))

			query.WriteString(` )`)

			if i < len(input.UpsertChannelStockInputs)-1 {
				query.WriteString(`, `)
			} else {
				query.WriteString(` `)
			}
		}

		query.WriteString(`
				ON CONFLICT ON CONSTRAINT channel_stock_pk
				DO
					UPDATE SET
						stock = EXCLUDED.stock,
						version = channel_stock.version + 1,
						updated_at = NOW()
			),
		`)

		query.WriteString(`
			delete_related_channel_stocks AS (
				DELETE FROM "temochi".channel_stock
				WHERE
		`)

		args = append(args, input.WarehouseID)
		query.WriteString(`warehouse_id = $` + strconv.Itoa(len(args)) + ` `)

		query.WriteString(`AND sku IN (`)
		for i := range input.UpsertInventoryInputs {
			args = append(args, input.UpsertChannelStockInputs[i].SKU)
			query.WriteString(`$` + strconv.Itoa(len(args)))
			if i < len(input.UpsertInventoryInputs)-1 {
				query.WriteString(`, `)
			} else {
				query.WriteString(`) `)
			}
		}

		query.WriteString(`AND CONCAT (sku, '#', gate_id, '#', channel_id) NOT IN (`)
		for i := range input.UpsertChannelStockInputs {
			args = append(args, concatenateSKUGateIDChannelID(input.UpsertChannelStockInputs[i]))
			query.WriteString(`$` + strconv.Itoa(len(args)))
			if i < len(input.UpsertChannelStockInputs)-1 {
				query.WriteString(`, `)
			} else {
				query.WriteString(`) `)
			}
		}
		query.WriteString(`
			)
		`)

		query.WriteString(`
			INSERT INTO "temochi".inventory (
				warehouse_id,
				sku,
				stock,
				buffer_stock
			)
			VALUES
		`)

		for i := range input.UpsertInventoryInputs {
			query.WriteString(`( `)

			args = append(args, input.WarehouseID)
			query.WriteString(`$` + strconv.Itoa(len(args)) + `, `)

			args = append(args, input.UpsertInventoryInputs[i].SKU)
			query.WriteString(`$` + strconv.Itoa(len(args)) + `, `)

			args = append(args, input.UpsertInventoryInputs[i].Stock)
			query.WriteString(`$` + strconv.Itoa(len(args)) + `, `)

			args = append(args, input.UpsertInventoryInputs[i].BufferStock)
			query.WriteString(`$` + strconv.Itoa(len(args)))

			query.WriteString(` )`)

			if i < len(input.UpsertInventoryInputs)-1 {
				query.WriteString(`, `)
			} else {
				query.WriteString(` `)
			}
		}

		query.WriteString(`
			ON CONFLICT ON CONSTRAINT inventory_pk
			DO
				UPDATE SET
					stock = EXCLUDED.stock,
					buffer_stock = EXCLUDED.buffer_stock,
					version = inventory.version + 1,
					updated_at = NOW()
		`)
	}

	return whitespaceNormalizer.ReplaceAllString(query.String(), " "), args
}

func concatenateSKUGateIDChannelID(item service.UpsertChannelStockInput) string {
	return item.SKU + "#" + item.GateID + "#" + item.ChannelID
}
