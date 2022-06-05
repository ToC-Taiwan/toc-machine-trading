BEGIN;

CREATE TABLE basic_calendar (
    "date" TIMESTAMP WITH Time Zone PRIMARY KEY,
    "is_trade_day" BOOLEAN NOT NULL
);

COMMIT;
