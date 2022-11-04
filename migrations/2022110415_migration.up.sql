CREATE TABLE history_tick_future (
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

CREATE INDEX history_tick_future_code_index ON history_tick_future USING btree ("code");

ALTER TABLE history_tick_future
ADD CONSTRAINT "fk_history_tick_future_future" FOREIGN KEY ("code") REFERENCES basic_future ("code");
