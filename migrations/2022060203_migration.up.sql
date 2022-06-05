BEGIN;

CREATE TABLE basic_stock (
    "number" VARCHAR PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "exchange" VARCHAR NOT NULL,
    "category" VARCHAR NOT NULL,
    "day_trade" BOOLEAN NOT NULL,
    "last_close" DECIMAL NOT NULL
);

COMMIT;
