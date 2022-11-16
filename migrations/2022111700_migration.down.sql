BEGIN;

ALTER TABLE trade_stock_balance DROP COLUMN IF EXISTS "manual";

ALTER TABLE trade_future_balance DROP COLUMN IF EXISTS "manual";

COMMIT;
