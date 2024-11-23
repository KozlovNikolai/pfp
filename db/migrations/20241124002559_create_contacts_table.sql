-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "contacts"(
    "id"  serial NOT NULL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "chat_id" INTEGER NOT NULL,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);

ALTER TABLE "contacts"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "contacts"
ADD FOREIGN KEY("chat_id") REFERENCES "chats"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "contacts" CASCADE;
-- +goose StatementEnd