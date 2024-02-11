BEGIN;

CREATE TABLE
    history_future_tick (
        "id" SERIAL PRIMARY KEY,
        "code" VARCHAR NOT NULL,
        "tick_time" TIMESTAMPTZ NOT NULL,
        "close" DECIMAL NOT NULL,
        "tick_type" INT NOT NULL,
        "volume" INT NOT NULL,
        "bid_price" DECIMAL NOT NULL,
        "bid_volume" INT NOT NULL,
        "ask_price" DECIMAL NOT NULL,
        "ask_volume" INT NOT NULL
    );

CREATE INDEX history_future_tick_code_index ON history_future_tick USING btree ("code");

ALTER TABLE history_future_tick ADD CONSTRAINT "fk_history_future_tick_future" FOREIGN KEY ("code") REFERENCES basic_future ("code");

COMMIT;
