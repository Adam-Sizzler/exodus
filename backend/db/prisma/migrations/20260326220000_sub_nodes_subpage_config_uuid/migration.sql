ALTER TABLE "sub_nodes"
    ADD COLUMN IF NOT EXISTS "subpage_config_uuid" UUID;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'sub_nodes_subpage_config_uuid_fkey'
    ) THEN
        ALTER TABLE "sub_nodes"
            ADD CONSTRAINT "sub_nodes_subpage_config_uuid_fkey"
            FOREIGN KEY ("subpage_config_uuid")
            REFERENCES "subscription_page_config"("uuid")
            ON DELETE SET NULL ON UPDATE CASCADE;
    END IF;
END$$;

CREATE INDEX IF NOT EXISTS "sub_nodes_subpage_config_uuid_idx"
    ON "sub_nodes"("subpage_config_uuid");
