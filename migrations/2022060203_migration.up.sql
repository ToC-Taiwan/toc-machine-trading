BEGIN;

CREATE TABLE basic_stock (
    "id" SERIAL PRIMARY KEY,
    "number" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    "exchange" VARCHAR NOT NULL,
    "category" VARCHAR NOT NULL,
    "day_trade" BOOLEAN NOT NULL,
    "last_close" DECIMAL NOT NULL
);

COMMIT;
