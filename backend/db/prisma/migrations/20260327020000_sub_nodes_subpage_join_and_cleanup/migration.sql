CREATE TABLE IF NOT EXISTS "sub_nodes_to_subscription_page_config" (
    "node_uuid" UUID NOT NULL,
    "subpage_config_uuid" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT "sub_nodes_to_subscription_page_config_pkey"
        PRIMARY KEY ("node_uuid", "subpage_config_uuid"),
    CONSTRAINT "sub_nodes_to_subscription_page_config_node_uuid_fkey"
        FOREIGN KEY ("node_uuid") REFERENCES "sub_nodes"("uuid")
            ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT "sub_nodes_to_subscription_page_config_subpage_config_uuid_fkey"
        FOREIGN KEY ("subpage_config_uuid") REFERENCES "subscription_page_config"("uuid")
            ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS "sub_nodes_to_subscription_page_config_node_uuid_key"
    ON "sub_nodes_to_subscription_page_config"("node_uuid");

INSERT INTO "sub_nodes_to_subscription_page_config" ("node_uuid", "subpage_config_uuid")
SELECT "uuid", "subpage_config_uuid"
FROM "sub_nodes"
WHERE "subpage_config_uuid" IS NOT NULL
ON CONFLICT ("node_uuid") DO UPDATE
SET "subpage_config_uuid" = EXCLUDED."subpage_config_uuid",
    "updated_at" = now();

DROP INDEX IF EXISTS "sub_nodes_subpage_config_uuid_idx";

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'sub_nodes_subpage_config_uuid_fkey'
    ) THEN
        ALTER TABLE "sub_nodes"
            DROP CONSTRAINT "sub_nodes_subpage_config_uuid_fkey";
    END IF;
END
$$;

ALTER TABLE IF EXISTS "sub_nodes"
    DROP COLUMN IF EXISTS "subpage_config_uuid";

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'sub_nodes_address_key'
    ) THEN
        ALTER TABLE "sub_nodes"
            DROP CONSTRAINT "sub_nodes_address_key";
    END IF;
END
$$;

DROP INDEX IF EXISTS "sub_nodes_address_key";

CREATE UNIQUE INDEX IF NOT EXISTS "sub_nodes_address_port_api_path_key"
    ON "sub_nodes"("address", "port", "api_path");

DROP TABLE IF EXISTS "config_profile_inbounds_to_subscription_connections";
DROP TABLE IF EXISTS "subscription_connections_traffic_usage_history";
DROP TABLE IF EXISTS "subscription_connections";
