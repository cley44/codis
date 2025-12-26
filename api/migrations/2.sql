CREATE TABLE "public"."workflow" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  "guild_id" text UNIQUE NOT NULL,
  "starting_discord_events" text[] NOT NULL,

  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "deleted_at" timestamp
);
