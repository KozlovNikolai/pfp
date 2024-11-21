-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "users"(
    "id"  serial NOT NULL PRIMARY KEY,
    "user_ext_id" text NOT NULL UNIQUE,
    "name" TEXT,
    "surname" TEXT,
    "email" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "users" CASCADE;
-- +goose StatementEnd