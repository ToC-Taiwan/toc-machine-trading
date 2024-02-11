BEGIN;

CREATE TABLE
    account_balance (
        "id" SERIAL PRIMARY KEY,
        "date" TIMESTAMPTZ NOT NULL,
        "balance" DECIMAL NOT NULL,
        "today_margin" DECIMAL NOT NULL,
        "available_margin" DECIMAL NOT NULL,
        "yesterday_margin" DECIMAL NOT NULL,
        "risk_indicator" DECIMAL NOT NULL
    );

CREATE TABLE
    account_settlement (
        "date" TIMESTAMPTZ PRIMARY KEY,
        "settlement" DECIMAL NOT NULL
    );

COMMIT;
