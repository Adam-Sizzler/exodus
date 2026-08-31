-- AlterTable config_profiles
ALTER TABLE "config_profiles" ADD COLUMN IF NOT EXISTS "tags" text[] NOT NULL DEFAULT '{}';

-- AlterTable internal_squads
ALTER TABLE "internal_squads" ADD COLUMN IF NOT EXISTS "tags" text[] NOT NULL DEFAULT '{}';

-- AlterTable external_squads
ALTER TABLE "external_squads" ADD COLUMN IF NOT EXISTS "tags" text[] NOT NULL DEFAULT '{}';

-- AlterTable node_plugin
ALTER TABLE "node_plugin" ADD COLUMN IF NOT EXISTS "tags" text[] NOT NULL DEFAULT '{}';

-- AlterTable subscription_templates
ALTER TABLE "subscription_templates" ADD COLUMN IF NOT EXISTS "tags" text[] NOT NULL DEFAULT '{}';

-- AlterTable subscription_page_config
ALTER TABLE "subscription_page_config" ADD COLUMN IF NOT EXISTS "tags" text[] NOT NULL DEFAULT '{}';

-- AlterTable hosts: add internal_squads_mode
ALTER TABLE "hosts" ADD COLUMN IF NOT EXISTS "internal_squads_mode" text NOT NULL DEFAULT 'EXCLUDE';

-- Rename internal_squad_host_exclusions to internal_squad_host_links if exists
DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables
        WHERE schemaname = 'public' AND tablename = 'internal_squad_host_exclusions'
    ) THEN
        ALTER TABLE "internal_squad_host_exclusions" RENAME TO "internal_squad_host_links";
    END IF;
END $$;

-- Create internal_squad_host_links if not exists
CREATE TABLE IF NOT EXISTS "internal_squad_host_links" (
    "host_uuid" UUID NOT NULL,
    "squad_uuid" UUID NOT NULL,

    CONSTRAINT "internal_squad_host_links_pkey" PRIMARY KEY ("host_uuid","squad_uuid"),
    CONSTRAINT "internal_squad_host_links_host_uuid_fkey" FOREIGN KEY ("host_uuid") REFERENCES "hosts"("uuid") ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT "internal_squad_host_links_squad_uuid_fkey" FOREIGN KEY ("squad_uuid") REFERENCES "internal_squads"("uuid") ON DELETE CASCADE ON UPDATE CASCADE
);
