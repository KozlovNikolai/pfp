-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "channel_twitter_profiles" (
    "id"  serial NOT NULL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"profile_id" VARCHAR(255) NOT NULL,
	"twitter_access_token" VARCHAR(255) NOT NULL,
	"twitter_access_token_secret" VARCHAR(255) NOT NULL,
	"account_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS "contacts" (
    "id"  serial NOT NULL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"email" VARCHAR(255) NOT NULL,
	"phone_number" VARCHAR(255) NOT NULL,
	"account_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"pubsub_token" VARCHAR(255) NOT NULL,
	"additional_atributes" INTEGER NOT NULL,
	"identifier" VARCHAR(255) NOT NULL,
	"source_id" BIGINT NOT NULL,
	"additional_attributes" JSON NOT NULL
);

CREATE INDEX "contacts_account_id_index"
ON "contacts" ("account_id");

CREATE TABLE IF NOT EXISTS "canned_responses" (
    "id"  serial NOT NULL PRIMARY KEY,
	"account_id" INTEGER NOT NULL,
	"short_code" VARCHAR(255) NOT NULL,
	"content" TEXT NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS "taggings" (
    "id"  serial NOT NULL PRIMARY KEY,
	"tag_id" INTEGER NOT NULL,
	"taggable_type" VARCHAR(255) NOT NULL,
	"taggable_id" INTEGER NOT NULL,
	"tagger_type" VARCHAR(255) NOT NULL,
	"tagger_id" INTEGER NOT NULL,
	"context" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL
);

CREATE INDEX "taggings_taggable_type_taggable_id_index"
ON "taggings" ("taggable_type", "taggable_id");

CREATE INDEX "taggings_tagger_type_tagger_id_index"
ON "taggings" ("tagger_type", "tagger_id");

CREATE INDEX "taggings_taggable_type_taggable_id_context_index"
ON "taggings" ("taggable_type", "taggable_id", "context");

CREATE INDEX "taggings_taggable_type_taggable_id_tagger_id_context_index"
ON "taggings" ("taggable_type", "taggable_id", "tagger_id", "context");

CREATE INDEX "taggings_tag_id_index"
ON "taggings" ("tag_id");

CREATE INDEX "taggings_taggable_type_index"
ON "taggings" ("taggable_type");

CREATE INDEX "taggings_taggable_id_index"
ON "taggings" ("taggable_id");

CREATE INDEX "taggings_tagger_id_index"
ON "taggings" ("tagger_id");

CREATE INDEX "taggings_context_index"
ON "taggings" ("context");

CREATE TABLE IF NOT EXISTS "subscriptions" (
    "id"  serial NOT NULL PRIMARY KEY,
	"pricing_version" VARCHAR(255) NOT NULL,
	"account_id" INTEGER NOT NULL,
	"expiry" TIMESTAMP NOT NULL,
	"billing_plan" VARCHAR(255) NOT NULL,
	"stripe_customer_id" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"state" INTEGER NOT NULL,
	"payment_source_added" SMALLINT NOT NULL
);


CREATE TABLE IF NOT EXISTS "conversations" (
    "id"  serial NOT NULL PRIMARY KEY,
	"account_id" INTEGER NOT NULL,
	"inbox_id" INTEGER NOT NULL,
	"status" INTEGER NOT NULL,
	"assignee_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"contact_id" BIGINT NOT NULL,
	"display_id" INTEGER NOT NULL,
	"user_last_seen_at" TIMESTAMP NOT NULL,
	"agent_last_seen_at" TIMESTAMP NOT NULL,
	"locked" SMALLINT NOT NULL,
	"contact_inbox_id" INTEGER NOT NULL,
	"additional_attributes" JSON NOT NULL
);

CREATE INDEX "conversations_account_id_index"
ON "conversations" ("account_id");

CREATE INDEX "conversations_contact_inbox_id_index"
ON "conversations" ("contact_inbox_id");

CREATE TABLE IF NOT EXISTS "contact_inboxes" (
    "id"  serial NOT NULL PRIMARY KEY,
	"contact_id" INTEGER NOT NULL,
	"inbox_id" INTEGER NOT NULL,
	"source_id" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);

CREATE INDEX "contact_inboxes_source_id_index"
ON "contact_inboxes" ("source_id");

CREATE TABLE IF NOT EXISTS "active_storage_blobs" (
    "id"  serial NOT NULL PRIMARY KEY,
	"key" VARCHAR(255) NOT NULL,
	"filename" VARCHAR(255) NOT NULL,
	"content_type" VARCHAR(255) NOT NULL,
	"metadata" TEXT NOT NULL,
	"byte_size" BIGINT NOT NULL,
	"checksum" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS "messages" (
    "id"  serial NOT NULL PRIMARY KEY,
	"content" TEXT NOT NULL,
	"account_id" INTEGER NOT NULL,
	"inbox_id" INTEGER NOT NULL,
	"conversation_id" INTEGER NOT NULL,
	"message_type" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"private" SMALLINT NOT NULL,
	"user_id" INTEGER NOT NULL,
	"status" INTEGER NOT NULL,
	"source_id" VARCHAR(255) NOT NULL,
	"content_type" INTEGER NOT NULL,
	"content_attributes" JSON NOT NULL,
	"contact_id" INTEGER NOT NULL
);

CREATE INDEX "messages_account_id_index"
ON "messages" ("account_id");

CREATE INDEX "messages_inbox_id_index"
ON "messages" ("inbox_id");

CREATE INDEX "messages_conversation_id_index"
ON "messages" ("conversation_id");

CREATE INDEX "messages_user_id_index"
ON "messages" ("user_id");

CREATE INDEX "messages_source_id_index"
ON "messages" ("source_id");

CREATE INDEX "messages_contact_id_index"
ON "messages" ("contact_id");

CREATE TABLE IF NOT EXISTS "channel_twilio_sms" (
    "id"  serial NOT NULL PRIMARY KEY,
	"phone_number" VARCHAR(255) NOT NULL,
	"auth_token" VARCHAR(255) NOT NULL,
	"account_sid" VARCHAR(255) NOT NULL,
	"account_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS "agent_bot_inboxes" (
    "id"  serial NOT NULL PRIMARY KEY,
	"inbox_id" INTEGER NOT NULL,
	"agent_bot_id" INTEGER NOT NULL,
	"status" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"account_id" INTEGER NOT NULL
);


CREATE TABLE IF NOT EXISTS "inboxes" (
    "id"  serial NOT NULL PRIMARY KEY,
	"channel_id" INTEGER NOT NULL,
	"account_id" INTEGER NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"channel_type" VARCHAR(255) NOT NULL,
	"enable_auto_assignment" SMALLINT NOT NULL
);

CREATE INDEX "inboxes_account_id_index"
ON "inboxes" ("account_id");

CREATE TABLE IF NOT EXISTS "events" (
    "id"  serial NOT NULL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"value" BYTEA NOT NULL,
	"account_id" INTEGER NOT NULL,
	"inbox_id" INTEGER NOT NULL,
	"user_id" INTEGER NOT NULL,
	"conversation_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);

CREATE INDEX "events_name_index"
ON "events" ("name");

CREATE INDEX "events_account_id_index"
ON "events" ("account_id");

CREATE INDEX "events_inbox_id_index"
ON "events" ("inbox_id");

CREATE INDEX "events_user_id_index"
ON "events" ("user_id");

CREATE INDEX "events_created_at_index"
ON "events" ("created_at");

CREATE TABLE IF NOT EXISTS "notification_settings" (
    "id"  serial NOT NULL PRIMARY KEY,
	"account_id" INTEGER NOT NULL,
	"user_id" INTEGER NOT NULL,
	"email_flags" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS "access_tokens" (
    "id"  serial NOT NULL PRIMARY KEY,
	"owner_type" VARCHAR(255) NOT NULL,
	"owner_id" BIGINT NOT NULL,
	"token" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);
   
CREATE INDEX "access_tokens_owner_type_owner_id_index"
ON "access_tokens" ("owner_type", "owner_id");

CREATE TABLE IF NOT EXISTS "accounts" (
    "id"  serial NOT NULL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"locale" INTEGER NOT NULL
);


CREATE TABLE IF NOT EXISTS "users" (
    "id"  serial NOT NULL PRIMARY KEY,
	"provider" VARCHAR(255) NOT NULL,
	"uid" VARCHAR(255) NOT NULL,
	"inviter" INTEGER NOT NULL,
	"encrypted_password" VARCHAR(255) NOT NULL,
	"reset_password_token" VARCHAR(255) NOT NULL,
	"reset_password_sent_at" TIMESTAMP NOT NULL,
	"remember_created_at" TIMESTAMP NOT NULL,
	"sign_in_count" INTEGER NOT NULL,
	"current_sign_in_at" TIMESTAMP NOT NULL,
	"last_sign_in_at" TIMESTAMP NOT NULL,
	"current_sign_in_ip" VARCHAR(255) NOT NULL,
	"last_sign_in_ip" VARCHAR(255) NOT NULL,
	"confirmation_token" VARCHAR(255) NOT NULL,
	"confirmed_at" TIMESTAMP NOT NULL,
	"confirmation_sent_at" TIMESTAMP NOT NULL,
	"unconfirmed_email" VARCHAR(255) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"nickname" VARCHAR(255) NOT NULL,
	"email" VARCHAR(255) NOT NULL,
	"tokens" JSON NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"pubsub_token" VARCHAR(255) NOT NULL
);

CREATE INDEX "users_email_index"
ON "users" ("email");

CREATE TABLE IF NOT EXISTS "inbox_members" (
    "id"  serial NOT NULL PRIMARY KEY,
	"user_id" INTEGER NOT NULL,
	"inbox_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);

CREATE INDEX "inbox_members_inbox_id_index"
ON "inbox_members" ("inbox_id");

CREATE TABLE IF NOT EXISTS "active_storage_attachments" (
    "id"  serial NOT NULL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"record_type" VARCHAR(255) NOT NULL,
	"record_id" BIGINT NOT NULL,
	"blob_id" BIGINT NOT NULL,
	"created_at" TIMESTAMP NOT NULL
);

CREATE INDEX "active_storage_attachments_blob_id_index"
ON "active_storage_attachments" ("blob_id");

CREATE TABLE IF NOT EXISTS "channel_web_widgets" (
    "id"  serial NOT NULL PRIMARY KEY,
	"website_name" VARCHAR(255) NOT NULL,
	"website_url" VARCHAR(255) NOT NULL,
	"account_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"website_token" VARCHAR(255) NOT NULL,
	"widget_color" VARCHAR(255) NOT NULL
);


CREATE TABLE IF NOT EXISTS "webhooks" (
    "id"  serial NOT NULL PRIMARY KEY,
	"account_id" INTEGER NOT NULL,
	"inbox_id" INTEGER NOT NULL,
	"url" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"webhook_type" INTEGER NOT NULL
);


CREATE TABLE IF NOT EXISTS "tags" (
    "id"  serial NOT NULL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"taggings_count" INTEGER NOT NULL
);


CREATE TABLE IF NOT EXISTS "agent_bots" (
    "id"  serial NOT NULL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"description" VARCHAR(255) NOT NULL,
	"outgoing_url" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS "attachments" (
    "id"  serial NOT NULL PRIMARY KEY,
	"file_type" INTEGER NOT NULL,
	"external_url" VARCHAR(255) NOT NULL,
	"coordinates_lat" BYTEA NOT NULL,
	"coordinates_long" BYTEA NOT NULL,
	"message_id" INTEGER NOT NULL,
	"account_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"fallback_title" VARCHAR(255) NOT NULL,
	"extension" VARCHAR(255) NOT NULL
);


CREATE TABLE IF NOT EXISTS "channel_facebook_pages" (
    "id"  serial NOT NULL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"page_id" VARCHAR(255) NOT NULL,
	"user_access_token" VARCHAR(255) NOT NULL,
	"page_access_token" VARCHAR(255) NOT NULL,
	"account_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);

CREATE INDEX "channel_facebook_pages_page_id_index"
ON "channel_facebook_pages" ("page_id");

CREATE TABLE IF NOT EXISTS "account_users" (
    "id"  serial NOT NULL PRIMARY KEY,
	"account_id" INTEGER NOT NULL,
	"user_id" INTEGER NOT NULL,
	"role" INTEGER NOT NULL,
	"inviter_id" BIGINT NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);

CREATE INDEX "account_users_account_id_index"
ON "account_users" ("account_id");

CREATE INDEX "account_users_user_id_index"
ON "account_users" ("user_id");

CREATE TABLE IF NOT EXISTS "telegram_bots" (
    "id"  serial NOT NULL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"auth_key" VARCHAR(255) NOT NULL,
	"account_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL
);


ALTER TABLE "subscriptions"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "conversations"
ADD FOREIGN KEY("inbox_id") REFERENCES "inboxes"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "telegram_bots"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "notification_settings"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "active_storage_attachments"
ADD FOREIGN KEY("blob_id") REFERENCES "active_storage_blobs"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "contacts"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "webhooks"
ADD FOREIGN KEY("inbox_id") REFERENCES "inboxes"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "attachments"
ADD FOREIGN KEY("message_id") REFERENCES "messages"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "contact_inboxes"
ADD FOREIGN KEY("inbox_id") REFERENCES "inboxes"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "agent_bot_inboxes"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "users"
ADD FOREIGN KEY("inviter") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "events"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "events"
ADD FOREIGN KEY("conversation_id") REFERENCES "conversations"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "channel_facebook_pages"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "webhooks"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "agent_bot_inboxes"
ADD FOREIGN KEY("agent_bot_id") REFERENCES "agent_bots"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "notification_settings"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "conversations"
ADD FOREIGN KEY("contact_inbox_id") REFERENCES "contact_inboxes"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "inbox_members"
ADD FOREIGN KEY("inbox_id") REFERENCES "inboxes"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "canned_responses"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "messages"
ADD FOREIGN KEY("conversation_id") REFERENCES "conversations"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "messages"
ADD FOREIGN KEY("inbox_id") REFERENCES "inboxes"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "account_users"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "conversations"
ADD FOREIGN KEY("contact_id") REFERENCES "contacts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "inbox_members"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "channel_twitter_profiles"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "account_users"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "conversations"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "agent_bot_inboxes"
ADD FOREIGN KEY("inbox_id") REFERENCES "inboxes"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "inboxes"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "messages"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "messages"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "attachments"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "events"
ADD FOREIGN KEY("inbox_id") REFERENCES "inboxes"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "channel_web_widgets"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "contact_inboxes"
ADD FOREIGN KEY("contact_id") REFERENCES "contacts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "channel_twilio_sms"
ADD FOREIGN KEY("account_id") REFERENCES "accounts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "messages"
ADD FOREIGN KEY("contact_id") REFERENCES "contacts"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "events"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "taggings"
ADD FOREIGN KEY("tag_id") REFERENCES "tags"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
-- +goose StatementEnd


-- ####################################################################################################
-- +goose Down
-- +goose StatementBegin
drop table IF EXISTS "telegram_bots" CASCADE;
drop table IF EXISTS "account_users" CASCADE;
drop table IF EXISTS "channel_facebook_pages" CASCADE;
drop table IF EXISTS "attachments" CASCADE;
drop table IF EXISTS "agent_bots" CASCADE;
drop table IF EXISTS "tags" CASCADE;
drop table IF EXISTS "webhooks" CASCADE;
drop table IF EXISTS "channel_web_widgets" CASCADE;
drop table IF EXISTS "active_storage_attachments" CASCADE;
drop table IF EXISTS "inbox_members" CASCADE;
drop table IF EXISTS "users" CASCADE;
drop table IF EXISTS "accounts" CASCADE;
drop table IF EXISTS "access_tokens" CASCADE;
drop table IF EXISTS "notification_settings" CASCADE;
drop table IF EXISTS "events" CASCADE;
drop table IF EXISTS "inboxes" CASCADE;
drop table IF EXISTS "agent_bot_inboxes" CASCADE;
drop table IF EXISTS "channel_twilio_sms" CASCADE;
drop table IF EXISTS "messages" CASCADE;
drop table IF EXISTS "active_storage_blobs" CASCADE;
drop table IF EXISTS "contact_inboxes" CASCADE;
drop table IF EXISTS "conversations" CASCADE;
drop table IF EXISTS "subscriptions" CASCADE;
drop table IF EXISTS "taggings" CASCADE;
drop table IF EXISTS "canned_responses" CASCADE;
drop table IF EXISTS "contacts" CASCADE;
drop table IF EXISTS "channel_twitter_profiles" CASCADE;
-- +goose StatementEnd