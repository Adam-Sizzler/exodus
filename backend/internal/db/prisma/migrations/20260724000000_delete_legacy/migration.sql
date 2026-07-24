-- Drop modules_settings table
DROP TABLE IF EXISTS "modules_settings";

-- Update exodus_settings
ALTER TABLE "exodus_settings" 
  DROP COLUMN IF EXISTS "tg_auth_settings",
  DROP COLUMN IF EXISTS "modules_settings";

-- Rename node_plugins to node_plugin
ALTER TABLE "node_plugins" RENAME TO "node_plugin";

-- Drop unique constraint on node_plugin.name
ALTER TABLE "node_plugin" DROP CONSTRAINT IF EXISTS "node_plugins_name_key";

-- Update users indices
DROP INDEX IF EXISTS "users_created_at_idx";
CREATE INDEX IF NOT EXISTS "users_expire_at_idx" ON "users"("expire_at");

-- Update admin table
ALTER TABLE "admin" DROP COLUMN IF EXISTS "session_ttl_minutes";

-- Update nodes table
ALTER TABLE "nodes" ALTER COLUMN "grpc_auth_token" DROP DEFAULT;
ALTER TABLE "nodes" ADD CONSTRAINT "nodes_active_plugin_uuid_fkey" FOREIGN KEY ("active_plugin_uuid") REFERENCES "node_plugin"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;

-- Update hosts table
ALTER TABLE "hosts" DROP COLUMN IF EXISTS "allow_insecure";

-- Update infra_billing_nodes
ALTER TABLE "infra_billing_nodes" ALTER COLUMN "node_uuid" DROP NOT NULL;
ALTER TABLE "infra_billing_nodes" ADD COLUMN IF NOT EXISTS "name" TEXT;

-- Update user_subscription_request_history indices
DROP INDEX IF EXISTS "user_subscription_request_history_request_at_idx";

-- Update user_meta and node_meta
ALTER TABLE "user_meta" ALTER COLUMN "metadata" DROP DEFAULT;
ALTER TABLE "node_meta" ALTER COLUMN "metadata" DROP DEFAULT;
