-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "unread_msgs"(
    "id"  serial NOT NULL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "message_id" INTEGER NOT NULL
);

ALTER TABLE "unread_msgs"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "unread_msgs"
ADD FOREIGN KEY("message_id") REFERENCES "messages"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "unread_msgs" CASCADE;
-- +goose StatementEnd