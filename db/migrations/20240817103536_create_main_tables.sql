-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "users"(
    "id"  serial NOT NULL PRIMARY KEY,
    "login" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "token" TEXT NOT NULL
);

ALTER TABLE users ADD CONSTRAINT unique_login UNIQUE (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table "users";
-- +goose StatementEnd