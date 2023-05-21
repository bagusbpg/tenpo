-- Deploy temochi:20230429_create_temochi_schema to pg

BEGIN;

CREATE SCHEMA IF NOT EXISTS temochi;

CREATE TABLE IF NOT EXISTS "temochi".inventory (
    warehouse_id    TEXT        NOT NULL,
    sku             TEXT        NOT NULL,
    stock           INTEGER     NOT NULL,
    buffer_stock    INTEGER     NOT NULL,
    version         BIGINT      DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT inventory_pk PRIMARY KEY (warehouse_id, sku),
    CONSTRAINT minimum_zero_stock CHECK (stock >= 0),
    CONSTRAINT valid_buffer_stock CHECK (buffer_stock <= stock)
);

CREATE TABLE IF NOT EXISTS "temochi".channel_stock (
    warehouse_id    TEXT        NOT NULL,
    sku             TEXT        NOT NULL,
    gate_id         TEXT        NOT NULL,
    channel_id      TEXT        NOT NULL,
    stock           INTEGER     NOT NULL,
    version         BIGINT      DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT channel_stock_pk PRIMARY KEY (warehouse_id, sku, gate_id, channel_id),
    FOREIGN KEY (warehouse_id, sku) REFERENCES "temochi".inventory (warehouse_id, sku) ON DELETE CASCADE
);

GRANT USAGE ON SCHEMA "temochi" TO "temochiapp";
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA "temochi" TO "temochiapp";

COMMIT;
