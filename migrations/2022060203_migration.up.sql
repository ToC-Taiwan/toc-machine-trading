BEGIN;

CREATE TABLE
    basic_calendar (
        "date" TIMESTAMPTZ PRIMARY KEY,
        "is_trade_day" BOOLEAN NOT NULL
    );

CREATE TABLE
    basic_stock (
        "number" VARCHAR PRIMARY KEY,
        "name" VARCHAR NOT NULL,
        "exchange" VARCHAR NOT NULL,
        "category" VARCHAR NOT NULL,
        "day_trade" BOOLEAN NOT NULL,
        "last_close" DECIMAL NOT NULL,
        "update_date" TIMESTAMPTZ NOT NULL
    );

CREATE TABLE
    basic_targets (
        "id" SERIAL PRIMARY KEY,
        "rank" INT NOT NULL,
        "stock_num" VARCHAR NOT NULL,
        "volume" INT NOT NULL,
        "trade_day" TIMESTAMPTZ
    );

CREATE INDEX basic_targets_trade_day_index ON basic_targets USING btree ("trade_day");

ALTER TABLE basic_targets ADD CONSTRAINT "fk_basic_targets_stock" FOREIGN KEY ("stock_num") REFERENCES basic_stock ("number");

CREATE TABLE
    history_stock_analyze (
        "id" SERIAL PRIMARY KEY,
        "stock_num" VARCHAR NOT NULL,
        "date" TIMESTAMPTZ,
        "quater_ma" DECIMAL NOT NULL
    );

CREATE INDEX history_stock_analyze_stock_num_index ON history_stock_analyze USING btree ("stock_num");

ALTER TABLE history_stock_analyze ADD CONSTRAINT "fk_history_stock_analyze_stock" FOREIGN KEY ("stock_num") REFERENCES basic_stock ("number");

CREATE TABLE
    history_stock_close (
        "id" SERIAL PRIMARY KEY,
        "date" TIMESTAMPTZ NOT NULL,
        "stock_num" VARCHAR NOT NULL,
        "close" DECIMAL NOT NULL
    );

CREATE INDEX history_stock_close_stock_num_index ON history_stock_close USING btree ("stock_num");

ALTER TABLE history_stock_close ADD CONSTRAINT "fk_history_stock_close_stock" FOREIGN KEY ("stock_num") REFERENCES basic_stock ("number");

CREATE TABLE
    history_stock_kbar (
        "id" SERIAL PRIMARY KEY,
        "stock_num" VARCHAR NOT NULL,
        "kbar_time" TIMESTAMPTZ NOT NULL,
        "open" DECIMAL NOT NULL,
        "high" DECIMAL NOT NULL,
        "low" DECIMAL NOT NULL,
        "close" DECIMAL NOT NULL,
        "volume" INT NOT NULL
    );

CREATE INDEX history_stock_kbar_stock_num_index ON history_stock_kbar USING btree ("stock_num");

ALTER TABLE history_stock_kbar ADD CONSTRAINT "fk_history_stock_kbar_stock" FOREIGN KEY ("stock_num") REFERENCES basic_stock ("number");

CREATE TABLE
    history_stock_tick (
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

CREATE INDEX history_stock_tick_stock_num_index ON history_stock_tick USING btree ("stock_num");

ALTER TABLE history_stock_tick ADD CONSTRAINT "fk_history_stock_tick_stock" FOREIGN KEY ("stock_num") REFERENCES basic_stock ("number");

CREATE TABLE
    trade_stock_order (
        "order_id" VARCHAR PRIMARY KEY,
        "status" INT NOT NULL,
        "order_time" TIMESTAMPTZ NOT NULL,
        "stock_num" VARCHAR NOT NULL,
        "action" INT NOT NULL,
        "price" DECIMAL NOT NULL,
        "lot" INT NOT NULL,
        "share" INT NOT NULL
    );

CREATE INDEX trade_stock_order_order_time_index ON trade_stock_order USING btree ("order_time");

ALTER TABLE trade_stock_order ADD CONSTRAINT "fk_trade_stock_order_stock" FOREIGN KEY ("stock_num") REFERENCES basic_stock ("number");

CREATE TABLE
    sinopac_event (
        "id" SERIAL PRIMARY KEY,
        "event_code" INT NOT NULL,
        "response" INT NOT NULL,
        "event" VARCHAR NOT NULL,
        "info" VARCHAR NOT NULL,
        "event_time" TIMESTAMPTZ NOT NULL
    );

CREATE TABLE
    trade_stock_balance (
        "id" SERIAL PRIMARY KEY,
        "trade_count" INT NOT NULL,
        "forward" INT NOT NULL,
        "reverse" INT NOT NULL,
        "original_balance" INT NOT NULL,
        "discount" INT NOT NULL,
        "total" INT NOT NULL,
        "trade_day" TIMESTAMPTZ NOT NULL
    );

CREATE INDEX trade_stock_balance_trade_day_index ON trade_stock_balance USING btree ("trade_day");

COMMIT;
