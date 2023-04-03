BEGIN;
CREATE TABLE history_future_close (
    "id" SERIAL PRIMARY KEY,
    "date" TIMESTAMPTZ NOT NULL,
    "code" VARCHAR NOT NULL,
    "close" DECIMAL NOT NULL
);
CREATE INDEX history_future_close_code_index ON history_future_close USING btree ("code");
ALTER TABLE history_future_close
ADD CONSTRAINT "fk_history_future_close_future" FOREIGN KEY ("code") REFERENCES basic_future ("code");
COMMIT;
