BEGIN;

CREATE TABLE basic_future (
    "code" VARCHAR PRIMARY KEY,
    "symbol" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    "category" VARCHAR NOT NULL,
    "delivery_month" VARCHAR NOT NULL,
    "delivery_date" VARCHAR NOT NULL,
    "underlying_kind" VARCHAR NOT NULL,
    "unit" INT NOT NULL,
    "limit_up" DECIMAL NOT NULL,
    "limit_down" DECIMAL NOT NULL,
    "reference" DECIMAL NOT NULL,
    "update_date" TIMESTAMPTZ NOT NULL
);

COMMIT;
