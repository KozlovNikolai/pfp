-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "unread_msgs"(
    "id"  serial NOT NULL PRIMARY KEY,
    "user_id" NUMERIC NOT NULL,
    "message_id" NUMERIC NOT NULL
);

-- ALTER TABLE myusers ADD CONSTRAINT unique_login UNIQUE (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "unread_msgs" CASCADE;
-- +goose StatementEnd