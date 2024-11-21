-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "chats"(
    "id"  serial NOT NULL PRIMARY KEY,
    "owner_id" INTEGER NOT NULL,
    "name" TEXT,
    "type" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "chats"
ADD FOREIGN KEY("owner_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if EXISTS "chats" CASCADE;
-- +goose StatementEnd