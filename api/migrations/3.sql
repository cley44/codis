CREATE TABLE "public"."node" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  "workflow_id" uuid NOT NULL,
  "type" text NOT NULL,
  "next_node_id" uuid,

  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "deleted_at" timestamp
);
