CREATE TABLE "public"."workflow" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  "starting_nodes_ids" uuid[] NOT NULL,
  "guild_id" text UNIQUE NOT NULL,
  "starting_discord_event" text NOT NULL,

  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "deleted_at" timestamp
);
