CREATE TABLE "public"."node" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  "workflow_id" uuid REFERENCES "public"."workflow"("id") NOT NULL,
  "type" text NOT NULL,
  "next_node_id" uuid REFERENCES "public"."node"("id"),

  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "deleted_at" timestamp
);

CREATE TABLE "public"."workflow_starting_node" (
  workflow_id uuid NOT NULL REFERENCES public.workflow(id) ON DELETE CASCADE,
  node_id uuid NOT NULL REFERENCES public.node(id) ON DELETE CASCADE,
  PRIMARY KEY (workflow_id, node_id)
);