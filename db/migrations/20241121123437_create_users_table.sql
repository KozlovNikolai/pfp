-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "users"(
    "id"  serial NOT NULL PRIMARY KEY,
    "user_ext_id" TEXT not NULL,
    "login" text not NULL,
    "password" TEXT,
    "account" text not NULL,
    "token" TEXT,
    "name" TEXT,
    "surname" TEXT,
    "email" TEXT NOT NULL,
    "user_type" TEXT NOT NULL,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);

CREATE unique INDEX "users_user_ext_id_account_index"
ON "users" ("user_ext_id", "account");

CREATE unique INDEX "users_login_account_index"
ON "users" ("login", "account");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "users" CASCADE;
-- +goose StatementEnd