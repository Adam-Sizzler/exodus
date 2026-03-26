ALTER TABLE "sub_nodes"
    ADD COLUMN IF NOT EXISTS "panel_api_token_uuid" UUID;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'sub_nodes_panel_api_token_uuid_fkey'
    ) THEN
        ALTER TABLE "sub_nodes"
            ADD CONSTRAINT "sub_nodes_panel_api_token_uuid_fkey"
            FOREIGN KEY ("panel_api_token_uuid")
            REFERENCES "api_tokens"("uuid")
            ON DELETE SET NULL ON UPDATE CASCADE;
    END IF;
END$$;

CREATE INDEX IF NOT EXISTS "sub_nodes_panel_api_token_uuid_idx"
    ON "sub_nodes"("panel_api_token_uuid");
