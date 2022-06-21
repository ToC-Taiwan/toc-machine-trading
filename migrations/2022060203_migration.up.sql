BEGIN;

CREATE TABLE basic_calendar (
    "date" TIMESTAMPTZ PRIMARY KEY,
    "is_trade_day" BOOLEAN NOT NULL
);

CREATE TABLE basic_stock (
    "number" VARCHAR PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "exchange" VARCHAR NOT NULL,
    "category" VARCHAR NOT NULL,
    "day_trade" BOOLEAN NOT NULL,
    "last_close" DECIMAL NOT NULL
);

CREATE TABLE basic_targets (
    "id" SERIAL PRIMARY KEY,
    "rank" INT NOT NULL,
    "stock_num" VARCHAR NOT NULL,
    "volume" INT NOT NULL,
    "subscribe" BOOLEAN NOT NULL,
    "real_time_add" BOOLEAN NOT NULL,
    "trade_day" TIMESTAMPTZ
);

CREATE TABLE history_close (
    "id" SERIAL PRIMARY KEY,
    "date" TIMESTAMPTZ NOT NULL,
    "stock_num" VARCHAR NOT NULL,
    "close" DECIMAL NOT NULL
);

CREATE TABLE history_kbar (
    "id" SERIAL PRIMARY KEY,
    "stock_num" VARCHAR NOT NULL,
    "kbar_time" TIMESTAMPTZ NOT NULL,
    "open" DECIMAL NOT NULL,
    "high" DECIMAL NOT NULL,
    "low" DECIMAL NOT NULL,
    "close" DECIMAL NOT NULL,
    "volume" INT NOT NULL
);

CREATE TABLE history_tick (
    "id" SERIAL PRIMARY KEY,
    "stock_num" VARCHAR NOT NULL,
    "tick_time" TIMESTAMPTZ NOT NULL,
    "close" DECIMAL NOT NULL,
    "tick_type" INT NOT NULL,
    "volume" INT NOT NULL,
    "bid_price" DECIMAL NOT NULL,
    "bid_volume" INT NOT NULL,
    "ask_price" DECIMAL NOT NULL,
    "ask_volume" INT NOT NULL
);

CREATE TABLE order (
    "order_id" VARCHAR PRIMARY KEY,
    "stock_num" VARCHAR NOT NULL,
    "action" INT NOT NULL,
    "price" DECIMAL NOT NULL,
    "quantity" INT NOT NULL,
    "status" INT NOT NULL,
    "order_time" TIMESTAMPTZ NOT NULL
);

CREATE TABLE sinopac_event (
    "id" SERIAL PRIMARY KEY,
    "event_code" INT NOT NULL,
    "response" INT NOT NULL,
    "event" VARCHAR NOT NULL,
    "info" VARCHAR NOT NULL,
    "event_time" TIMESTAMPTZ NOT NULL
);

CREATE TABLE trade_balance (
    "id" SERIAL PRIMARY KEY,
    "trade_count" INT NOT NULL,
    "forward" INT NOT NULL,
    "reverse" INT NOT NULL,
    "original_balance" INT NOT NULL,
    "discount" INT NOT NULL,
    "total" INT NOT NULL,
    "trade_day" TIMESTAMPTZ NOT NULL
);

COMMIT;
