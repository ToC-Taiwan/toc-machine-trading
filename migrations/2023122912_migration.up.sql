BEGIN;
CREATE TABLE system_account (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR NOT NULL,
    "password" VARCHAR NOT NULL,
    "email" VARCHAR NOT NULL,
    "email_verified" BOOLEAN NOT NULL DEFAULT FALSE,
    "auth_trade" BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE TABLE system_push_token (
    "id" SERIAL PRIMARY KEY,
    "created" TIMESTAMPTZ NOT NULL,
    "token" VARCHAR NOT NULL UNIQUE,
    "user_id" INT NOT NULL
);
CREATE INDEX system_push_token_user_index ON system_push_token USING btree ("created");
ALTER TABLE system_push_token
ADD CONSTRAINT "fk_system_push_token_user" FOREIGN KEY ("user_id") REFERENCES system_account ("id");
COMMIT;
