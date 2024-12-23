-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "contacts"(
    "id"  serial NOT NULL PRIMARY KEY,
    "account_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "name" TEXT,
    "surname" TEXT,
    "phone" TEXT,
    "email" TEXT
);

CREATE unique INDEX "contacts_account_id_user_id_index"
ON "contacts" ("account_id", "user_id");

ALTER TABLE "contacts"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "contacts"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "contacts" CASCADE;
-- +goose StatementEnd