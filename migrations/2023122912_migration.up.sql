BEGIN;
CREATE TABLE system_account (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR NOT NULL,
    "password" VARCHAR NOT NULL,
    "email" VARCHAR NOT NULL,
    "push_token" VARCHAR,
    "email_verified" BOOLEAN NOT NULL DEFAULT FALSE,
    "auth_trade" BOOLEAN NOT NULL DEFAULT FALSE
);
COMMIT;
