-- Verify temochi:20230429_create_temochi_schema on pg

BEGIN;

DO $$DECLARE warehouse_id TEXT;
BEGIN
    INSERT INTO "temochi".inventory (warehouse_id, sku, stock, buffer_stock)
    VALUES ('dummy-warehouse-id', 'dummy-sku', 10, 0);

    warehouse_id = (
        SELECT
            inventory.warehouse_id
        FROM
            "temochi".inventory
        WHERE
            sku = 'dummy-sku'
    );

    ASSERT warehouse_id = 'dummy-warehouse-id';
END$$;

DO $$DECLARE warehouse_id TEXT;
BEGIN
    INSERT INTO "temochi".channel_stock (warehouse_id, sku, gate_id, channel_id, stock)
    VALUES ('dummy-warehouse-id', 'dummy-sku', 'dummy-gate-id', 'dummy-channel-id', 10);

    warehouse_id = (
        SELECT
            channel_stock.warehouse_id
        FROM
            "temochi".channel_stock
        WHERE
            sku = 'dummy-sku'
    );

    ASSERT warehouse_id = 'dummy-warehouse-id';
END$$;

ROLLBACK;
