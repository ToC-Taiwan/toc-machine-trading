BEGIN;

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
    "trade_day" TIMESTAMP WITH Time Zone,
    "rank" INT NOT NULL,
    "volume" INT NOT NULL,
    "subscribe" BOOLEAN NOT NULL,
    "real_time_add" BOOLEAN NOT NULL
);

CREATE TABLE basic_calendar (
    "date" TIMESTAMP WITH Time Zone PRIMARY KEY,
    "is_trade_day" BOOLEAN NOT NULL
);

CREATE TABLE sinopac_event (
    "id" SERIAL PRIMARY KEY,
    "event_code" INT NOT NULL,
    "response" INT NOT NULL,
    "event" VARCHAR NOT NULL,
    "info" VARCHAR NOT NULL
);

COMMIT;
