-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "chat_members"(
    "id"  serial NOT NULL PRIMARY KEY,
    "chat_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "role" TEXT,
    "last_read_msg_id" INTEGER,
    "notifications" BOOLEAN,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);

CREATE unique INDEX "chat_members_chat_id_user_id_index"
ON "chat_members" ("chat_id", "user_id");

ALTER TABLE "chat_members"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "chat_members"
ADD FOREIGN KEY("chat_id") REFERENCES "chats"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "chat_members" CASCADE;
-- +goose StatementEnd