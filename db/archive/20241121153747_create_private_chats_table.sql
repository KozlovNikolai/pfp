-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "private_chats"(
    "id"  serial NOT NULL PRIMARY KEY,
    "user1_id" INTEGER NOT NULL,
    "user2_id" INTEGER NOT NULL,
    "msg_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "private_chats" CASCADE;
-- +goose StatementEnd