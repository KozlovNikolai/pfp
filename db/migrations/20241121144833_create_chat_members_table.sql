-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "chat_members"(
    "id"  serial NOT NULL PRIMARY KEY,
    "chat_id" NUMERIC NOT NULL,
    "user_id" NUMERIC NOT NULL,
    "role" TEXT,
    "last_read_msg_id" NUMERIC,
    "notifications" BOOLEAN,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

-- ALTER TABLE myusers ADD CONSTRAINT unique_login UNIQUE (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "chat_members" CASCADE;
-- +goose StatementEnd