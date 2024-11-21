-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "notifications"(
    "id"  serial NOT NULL PRIMARY KEY,
    "user_id" NUMERIC NOT NULL,
    "chat_id" NUMERIC NOT NULL,
    "msg_id" NUMERIC NOT NULL,
    "text" TEXT,
    "is_deleted" BOOLEAN,
    "notification_type" TEXT NOT NULL,
    "is_seen" BOOLEAN,
    "created_at" TIMESTAMP NOT NULL
);

-- ALTER TABLE myusers ADD CONSTRAINT unique_login UNIQUE (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "notifications" CASCADE;
-- +goose StatementEnd