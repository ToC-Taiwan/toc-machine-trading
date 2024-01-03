BEGIN;
CREATE TABLE basic_future (
    "code" VARCHAR PRIMARY KEY,
    "symbol" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    "category" VARCHAR NOT NULL,
    "delivery_month" VARCHAR NOT NULL,
    "delivery_date" TIMESTAMPTZ NOT NULL,
    "underlying_kind" VARCHAR NOT NULL,
    "unit" INT NOT NULL,
    "limit_up" DECIMAL NOT NULL,
    "limit_down" DECIMAL NOT NULL,
    "reference" DECIMAL NOT NULL,
    "update_date" TIMESTAMPTZ NOT NULL
);
CREATE TABLE trade_future_order (
    "order_id" VARCHAR PRIMARY KEY,
    "status" INT NOT NULL,
    "order_time" TIMESTAMPTZ NOT NULL,
    "code" VARCHAR NOT NULL,
    "action" INT NOT NULL,
    "price" DECIMAL NOT NULL,
    "position" INT NOT NULL
);
CREATE INDEX trade_future_order_order_time_index ON trade_future_order USING btree ("order_time");
ALTER TABLE trade_future_order
ADD CONSTRAINT "fk_trade_future_order_future" FOREIGN KEY ("code") REFERENCES basic_future ("code");
CREATE TABLE trade_future_balance (
    "id" SERIAL PRIMARY KEY,
    "trade_count" INT NOT NULL,
    "forward" INT NOT NULL,
    "reverse" INT NOT NULL,
    "total" INT NOT NULL,
    "trade_day" TIMESTAMPTZ NOT NULL
);
CREATE INDEX trade_future_balance_trade_day_index ON trade_future_balance USING btree ("trade_day");
COMMIT;
