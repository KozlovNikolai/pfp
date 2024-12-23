-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "users"(
    "id"  serial NOT NULL PRIMARY KEY,
    "user_ext_id" INTEGER not NULL,
    "login" text not NULL,
    "password" TEXT,
    "profile" TEXT not NULL,
    "name" TEXT,
    "surname" TEXT,
    "email" TEXT NOT NULL,
    "user_type" TEXT NOT NULL,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);

-- CREATE unique INDEX "users_user_ext_id_account_index"
-- ON "users" ("user_ext_id", "account");

-- CREATE unique INDEX "users_login_account_index"
-- ON "users" ("login", "account");

CREATE unique INDEX "users_login_account_user_ext_id_index"
ON "users" ("login", "profile", "user_ext_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "users" CASCADE;
-- +goose StatementEnd