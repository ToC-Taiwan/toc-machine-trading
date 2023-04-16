BEGIN;
CREATE TABLE inventory_stock (
    "id" SERIAL PRIMARY KEY,
    "updated" TIMESTAMPTZ NOT NULL,
    "bank_id" INT NOT NULL,
    "avg_price" DECIMAL NOT NULL,
    "quantity" INT NOT NULL,
    "stock_num" VARCHAR NOT NULL
);
CREATE TABLE inventory_future (
    "id" SERIAL PRIMARY KEY,
    "updated" TIMESTAMPTZ NOT NULL,
    "bank_id" INT NOT NULL,
    "avg_price" DECIMAL NOT NULL,
    "quantity" INT NOT NULL,
    "code" VARCHAR NOT NULL
);
COMMIT;
