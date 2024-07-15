-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "id" DROP DEFAULT, ALTER COLUMN "id" TYPE uuid, ALTER COLUMN "id" SET DEFAULT gen_random_uuid();
-- Drop sequence used by serial column "id"
DROP SEQUENCE IF EXISTS "public"."users_id_seq";
