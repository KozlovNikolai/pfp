-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "chats"(
    "id"  serial NOT NULL PRIMARY KEY,
    "account_id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "chat_type" TEXT NOT NULL,
    "last_message_id" INTEGER,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);

ALTER TABLE "chats"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

CREATE unique INDEX "chats_name_chat_type_index"
ON "chats" ("name", "chat_type");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "chats" CASCADE;
-- +goose StatementEnd