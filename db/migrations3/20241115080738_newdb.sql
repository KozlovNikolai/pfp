-- +goose Up
-- +goose StatementBegin

CREATE TABLE "Chats" (
	"id" serial NOT NULL PRIMARY KEY,
	"uuid" uuid NOT NULL UNIQUE,
	"name" text NOT NULL,
	"created_at" timestamp NOT NULL,
	"update_at" timestamp NOT NULL
);

CREATE TABLE "Files" (
	"id" serial NOT NULL PRIMARY KEY,
	"uuid" uuid NOT NULL UNIQUE,
	"url" text NOT NULL,
	"file_type" int4 NOT NULL,
	"size" int4 NOT NULL,
	"uploaded_at" timestamp NOT NULL
);

CREATE TABLE "Users" (
	"id" serial  NOT NULL PRIMARY KEY,
	"uuid" uuid NOT NULL UNIQUE,
	"name" text NOT NULL,
	"email" text NOT NULL UNIQUE,
	"password_hash" text NOT NULL,
	"created_at" timestamp NOT NULL
);

CREATE TABLE "ChatMembers" (
	"id" serial NOT NULL PRIMARY KEY,
	"chat_id" uuid NOT NULL,
	"user_id" uuid NOT NULL,
	"role" text NOT NULL,
	"last_read_message_id" uuid NOT NULL,
	"notifications" bool NOT NULL,
	CONSTRAINT "ChatMembers_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES "Chats"("uuid"),
	CONSTRAINT "ChatMembers_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "Users"("uuid")
);


CREATE TABLE "Devices" (
	"id" serial NOT NULL PRIMARY KEY,
	"uuid" uuid NOT NULL UNIQUE,
	"user_id" uuid NOT NULL,
	"device_type" int NOT NULL,
	"device_token" text NOT NULL,
	"last_online" timestamp,
	"is_active" bool NOT NULL,
	CONSTRAINT "devices_users_fk" FOREIGN KEY ("user_id") REFERENCES "Users"("uuid")
);

CREATE TABLE "Messages" (
	"id" serial NOT NULL PRIMARY KEY,
	"uuid" uuid NOT NULL UNIQUE,
	"chat_id" uuid NOT NULL,
	"sender_id" uuid NOT NULL,
	"content" text NOT NULL,
	"created_at" timestamp NOT NULL,
	"is_deleted" bool,
	"file_id" uuid,
	"message_type" int NOT NULL,
	"seen_by" jsonb,
	"device_read_status" jsonb,
	CONSTRAINT "Messages_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES "Chats"("uuid"),
	CONSTRAINT "Messages_file_id_fkey" FOREIGN KEY ("file_id") REFERENCES "Files"("uuid"),
	CONSTRAINT "Messages_sender_id_fkey" FOREIGN KEY ("sender_id") REFERENCES "Users"("uuid")
);


CREATE TABLE "Notifications" (
	"id" serial NOT NULL PRIMARY KEY,
	"uuid" uuid NOT NULL UNIQUE,
	"user_id" uuid NOT NULL,
	"chat_id" uuid NOT NULL,
	"message_id" uuid NOT NULL,
	"created_at" timestamp NOT NULL,
	"notification_type" int,
	"is_seen" bool NULL,
	CONSTRAINT "Notifications_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES "Chats"("uuid"),
	CONSTRAINT "Notifications_message_id_fkey" FOREIGN KEY ("message_id") REFERENCES "Messages"("uuid"),
	CONSTRAINT "Notifications_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "Users"("uuid")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table IF EXISTS "Notifications" CASCADE;
drop table IF EXISTS "Messages" CASCADE;
drop table IF EXISTS "Devices" CASCADE;
drop table IF EXISTS "ChatMembers" CASCADE;
drop table IF EXISTS "Users" CASCADE;
drop table IF EXISTS "Files" CASCADE;
drop table IF EXISTS "Chats" CASCADE;

-- +goose StatementEnd
