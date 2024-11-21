-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "myusers"(
    "id"  serial NOT NULL PRIMARY KEY,
    "login" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "token" TEXT NOT NULL
);

ALTER TABLE myusers ADD CONSTRAINT unique_login UNIQUE (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table "myusers";
-- +goose StatementEnd