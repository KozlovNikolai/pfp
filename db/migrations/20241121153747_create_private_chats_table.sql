-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "private_chats"(
    "id"  serial NOT NULL PRIMARY KEY,
    "user1_id" NUMERIC NOT NULL,
    "user2_id" NUMERIC NOT NULL,
    "msg_id" NUMERIC NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

-- ALTER TABLE myusers ADD CONSTRAINT unique_login UNIQUE (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "private_chats" CASCADE;
-- +goose StatementEnd