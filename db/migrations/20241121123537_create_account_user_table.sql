-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "account_user"(
    "id"  serial NOT NULL PRIMARY KEY,
    "account_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "role" TEXT NOT NULL,
    "inviter_id" INTEGER NOT NULL,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);

CREATE unique INDEX "account_user_account_id_user_id_index"
ON "account_user" ("account_id", "user_id");

ALTER TABLE "account_user"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "account_user"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "account_user"
ADD FOREIGN KEY("inviter_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "account_user" CASCADE;
-- +goose StatementEnd