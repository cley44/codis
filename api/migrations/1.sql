-- Create "user" table
CREATE TABLE "public"."user" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  "username" text UNIQUE NOT NULL,
  "email" text UNIQUE NOT NULL,
  "display_username" text,
  "discord_id" text UNIQUE NOT NULL,
  "discord_avatar" text,
  "discord_session" jsonb,
  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "deleted_at" timestamp
);
