-- AlterTable nodes
ALTER TABLE "nodes" ADD COLUMN IF NOT EXISTS "ips" JSONB NOT NULL DEFAULT '[]';

-- AlterTable hosts
ALTER TABLE "hosts" ADD COLUMN IF NOT EXISTS "mapper" JSONB NOT NULL DEFAULT '{}';

-- AlterTable node_plugin
ALTER TABLE "node_plugin" ADD COLUMN IF NOT EXISTS "shared_lists" JSONB NOT NULL DEFAULT '[]';
