-- SRS lists: remove obsolete type field and add enable/disable flag.
ALTER TABLE "srs_lists"
    ADD COLUMN IF NOT EXISTS "is_enabled" BOOLEAN NOT NULL DEFAULT true;

ALTER TABLE "srs_lists"
    DROP COLUMN IF EXISTS "type";
