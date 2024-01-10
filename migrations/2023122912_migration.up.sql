BEGIN;
CREATE TABLE system_account (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR NOT NULL,
    "password" VARCHAR NOT NULL,
    "email" VARCHAR NOT NULL,
    "email_verified" BOOLEAN NOT NULL DEFAULT FALSE
);
COMMIT;
