BEGIN;
CREATE TABLE account_balance (
    "id" SERIAL PRIMARY KEY,
    "date" TIMESTAMPTZ NOT NULL,
    "balance" DECIMAL NOT NULL,
    "today_margin" DECIMAL NOT NULL,
    "available_margin" DECIMAL NOT NULL,
    "yesterday_margin" DECIMAL NOT NULL,
    "risk_indicator" DECIMAL NOT NULL,
    "bank_id" INT NOT NULL
);
CREATE TABLE account_settlement (
    "date" TIMESTAMPTZ PRIMARY KEY,
    "sinopac" DECIMAL NOT NULL,
    "fugle" DECIMAL NOT NULL
);
COMMIT;
