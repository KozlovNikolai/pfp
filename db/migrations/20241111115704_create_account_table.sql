-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "accounts"(
    "id"  serial NOT NULL PRIMARY KEY,
    "name" TEXT not NULL,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "accounts" CASCADE;
-- +goose StatementEnd