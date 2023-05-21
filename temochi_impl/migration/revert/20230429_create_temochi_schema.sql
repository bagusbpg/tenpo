-- Revert temochi:20230429_create_temochi_schema from pg

BEGIN;

DROP TABLE IF EXISTS "temochi".channel_stock;

DROP TABLE IF EXISTS "temochi".inventory;

DROP SCHEMA IF EXISTS temochi;

COMMIT;
