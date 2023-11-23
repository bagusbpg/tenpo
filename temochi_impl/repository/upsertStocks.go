package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type UpsertStockDBInput struct {
	WarehouseID              string
	UpsertInventoryInputs    []UpsertInventoryInput
	UpsertChannelStockInputs []UpsertChannelStockInput
}

type UpsertInventoryInput struct {
	SKU         string
	Stock       uint32
	BufferStock uint32
}

type UpsertChannelStockInput struct {
	SKU       string
	GateID    string
	ChannelID string
	Stock     uint32
}

type UpsertStockDBOutput struct{}

func (ths *repository) UpsertStock(ctx context.Context, input UpsertStockDBInput, output *UpsertStockDBOutput) error {
	query, args := buildUpsertStocksQuery(input)

	_, err := ths.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed executing UpsertInventory query: %s", err.Error())
	}

	return nil
}

func buildUpsertStocksQuery(input UpsertStockDBInput) (string, []interface{}) {
	query := strings.Builder{}
	args := make([]interface{}, 0)

	if len(input.UpsertChannelStockInputs) > 0 {
		query.WriteString(`
			WITH
				upsert_channel_stock AS (` + buildUpsertChannelStocksQuery(input, &args) + `),
				delete_related_channel_stocks AS (` + buildDeleteRelatedChannelStocksQuery(input, &args, false) + `)
		`)
	} else {
		query.WriteString(`
			WITH
				delete_all_channel_stocks AS (` + buildDeleteRelatedChannelStocksQuery(input, &args, true) + `)
		`)
	}

	query.WriteString(buildUpsertInventoryQuery(input, &args))

	return whitespaceNormalizer.ReplaceAllString(query.String(), " "), args
}

func buildUpsertChannelStocksQuery(input UpsertStockDBInput, args *[]interface{}) string {
	query := strings.Builder{}

	query.WriteString(`
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

		*args = append(*args, input.WarehouseID)
		query.WriteString(`$` + strconv.Itoa(len(*args)) + `, `)

		*args = append(*args, input.UpsertChannelStockInputs[i].SKU)
		query.WriteString(`$` + strconv.Itoa(len(*args)) + `, `)

		*args = append(*args, input.UpsertChannelStockInputs[i].GateID)
		query.WriteString(`$` + strconv.Itoa(len(*args)) + `, `)

		*args = append(*args, input.UpsertChannelStockInputs[i].ChannelID)
		query.WriteString(`$` + strconv.Itoa(len(*args)) + `, `)

		*args = append(*args, input.UpsertChannelStockInputs[i].Stock)
		query.WriteString(`$` + strconv.Itoa(len(*args)))

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
	`)

	return query.String()
}

func buildDeleteRelatedChannelStocksQuery(input UpsertStockDBInput, args *[]interface{}, deleteAll bool) string {
	query := strings.Builder{}

	query.WriteString(`
		DELETE FROM "temochi".channel_stock
		WHERE
	`)

	*args = append(*args, input.WarehouseID)
	query.WriteString(`warehouse_id = $` + strconv.Itoa(len(*args)) + ` `)

	query.WriteString(`AND sku IN (`)
	for i := range input.UpsertInventoryInputs {
		*args = append(*args, input.UpsertChannelStockInputs[i].SKU)
		query.WriteString(`$` + strconv.Itoa(len(*args)))
		if i < len(input.UpsertInventoryInputs)-1 {
			query.WriteString(`, `)
		} else {
			query.WriteString(`) `)
		}
	}

	if deleteAll {
		return query.String()
	}

	query.WriteString(`AND CONCAT (sku, '#', gate_id, '#', channel_id) NOT IN (`)
	for i := range input.UpsertChannelStockInputs {
		*args = append(*args, concatenateSKUGateIDChannelID(input.UpsertChannelStockInputs[i]))
		query.WriteString(`$` + strconv.Itoa(len(*args)))
		if i < len(input.UpsertChannelStockInputs)-1 {
			query.WriteString(`, `)
		} else {
			query.WriteString(`) `)
		}
	}

	return query.String()
}

func buildUpsertInventoryQuery(input UpsertStockDBInput, args *[]interface{}) string {
	query := strings.Builder{}

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

		*args = append(*args, input.WarehouseID)
		query.WriteString(`$` + strconv.Itoa(len(*args)) + `, `)

		*args = append(*args, input.UpsertInventoryInputs[i].SKU)
		query.WriteString(`$` + strconv.Itoa(len(*args)) + `, `)

		*args = append(*args, input.UpsertInventoryInputs[i].Stock)
		query.WriteString(`$` + strconv.Itoa(len(*args)) + `, `)

		*args = append(*args, input.UpsertInventoryInputs[i].BufferStock)
		query.WriteString(`$` + strconv.Itoa(len(*args)))

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

	return query.String()
}

func concatenateSKUGateIDChannelID(item UpsertChannelStockInput) string {
	return item.SKU + "#" + item.GateID + "#" + item.ChannelID
}
