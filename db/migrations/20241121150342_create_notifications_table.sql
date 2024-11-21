-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "notifications"(
    "id"  serial NOT NULL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "chat_id" INTEGER NOT NULL,
    "msg_id" INTEGER NOT NULL,
    "text" TEXT,
    "is_deleted" BOOLEAN,
    "notification_type" TEXT NOT NULL,
    "is_seen" BOOLEAN,
    "created_at" TIMESTAMP NOT NULL
);

ALTER TABLE "notifications"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "notifications"
ADD FOREIGN KEY("chat_id") REFERENCES "chats"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "notifications"
ADD FOREIGN KEY("msg_id") REFERENCES "messages"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "notifications" CASCADE;
-- +goose StatementEnd