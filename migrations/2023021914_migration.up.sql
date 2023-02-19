BEGIN;

CREATE TABLE basic_option (
    "code" VARCHAR PRIMARY KEY,
    "symbol" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    "category" VARCHAR NOT NULL,
    "delivery_month" VARCHAR NOT NULL,
    "delivery_date" TIMESTAMPTZ NOT NULL,
    "strike_price" DECIMAL NOT NULL,
    "option_right" VARCHAR NOT NULL,
    "underlying_kind" VARCHAR NOT NULL,
    "unit" INT NOT NULL,
    "limit_up" DECIMAL NOT NULL,
    "limit_down" DECIMAL NOT NULL,
    "reference" DECIMAL NOT NULL,
    "update_date" TIMESTAMPTZ NOT NULL
);

CREATE TABLE trade_option_order (
    "order_id" VARCHAR PRIMARY KEY,
    "status" INT NOT NULL,
    "order_time" TIMESTAMPTZ NOT NULL,
    "code" VARCHAR NOT NULL,
    "action" INT NOT NULL,
    "price" DECIMAL NOT NULL,
    "quantity" INT NOT NULL
);

CREATE INDEX trade_option_order_order_time_index ON trade_option_order USING btree ("order_time");

ALTER TABLE trade_option_order
ADD CONSTRAINT "fk_trade_option_order_option" FOREIGN KEY ("code") REFERENCES basic_option ("code");

CREATE TABLE trade_option_balance (
    "id" SERIAL PRIMARY KEY,
    "trade_count" INT NOT NULL,
    "forward" INT NOT NULL,
    "reverse" INT NOT NULL,
    "total" INT NOT NULL,
    "trade_day" TIMESTAMPTZ NOT NULL
);

CREATE INDEX trade_option_balance_trade_day_index ON trade_option_balance USING btree ("trade_day");

COMMIT;
