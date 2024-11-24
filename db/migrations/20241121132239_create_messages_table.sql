-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "messages"(
    "id"  serial NOT NULL PRIMARY KEY,
    "sender_id" INTEGER NOT NULL,
    "chat_id" INTEGER NOT NULL,
    "msg_type" TEXT NOT NULL,
    "text" TEXT,
    "is_deleted" BOOLEAN,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);

ALTER TABLE "messages"
ADD FOREIGN KEY("sender_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "messages"
ADD FOREIGN KEY("chat_id") REFERENCES "chats"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "messages" CASCADE;
-- +goose StatementEnd