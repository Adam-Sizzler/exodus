-- AlterTable
ALTER TABLE "srs_lists" DROP CONSTRAINT IF EXISTS "srs_lists_tag_key";
DROP INDEX IF EXISTS "srs_lists_tag_key";
ALTER TABLE "srs_lists" DROP COLUMN IF EXISTS "tag";

ALTER TABLE "srs_lists" ADD COLUMN IF NOT EXISTS "tags" text[] NOT NULL DEFAULT '{}';
