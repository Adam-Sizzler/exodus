DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'sub_nodes_panel_api_token_uuid_fkey'
    ) THEN
        ALTER TABLE "sub_nodes"
            DROP CONSTRAINT "sub_nodes_panel_api_token_uuid_fkey";
    END IF;
END$$;

DROP INDEX IF EXISTS "sub_nodes_panel_api_token_uuid_idx";

ALTER TABLE "sub_nodes"
    DROP COLUMN IF EXISTS "panel_api_token_uuid";
