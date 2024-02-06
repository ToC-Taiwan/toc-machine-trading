BEGIN;
CREATE TABLE inventory_stock (
    "id" SERIAL PRIMARY KEY,
    "updated" TIMESTAMPTZ NOT NULL,
    "avg_price" DECIMAL NOT NULL,
    "lot" INT NOT NULL,
    "share" INT NOT NULL,
    "stock_num" VARCHAR NOT NULL
);
CREATE INDEX inventory_stock_stock_num_index ON inventory_stock USING btree ("stock_num");
ALTER TABLE inventory_stock
ADD CONSTRAINT "fk_inventory_stock_future" FOREIGN KEY ("stock_num") REFERENCES basic_stock ("number");
CREATE TABLE inventory_future (
    "id" SERIAL PRIMARY KEY,
    "updated" TIMESTAMPTZ NOT NULL,
    "avg_price" DECIMAL NOT NULL,
    "position" INT NOT NULL,
    "code" VARCHAR NOT NULL
);
CREATE INDEX inventory_future_code_index ON inventory_future USING btree ("code");
ALTER TABLE inventory_future
ADD CONSTRAINT "fk_inventory_future_future" FOREIGN KEY ("code") REFERENCES basic_future ("code");
COMMIT;
