BEGIN;

ALTER TABLE trade_stock_order DROP COLUMN IF EXISTS "manual";

ALTER TABLE trade_future_order DROP COLUMN IF EXISTS "manual";

COMMIT;
