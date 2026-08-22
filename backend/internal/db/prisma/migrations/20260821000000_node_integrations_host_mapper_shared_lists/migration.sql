-- 1. Create integrations table
CREATE TABLE IF NOT EXISTS "integrations" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" VARCHAR(30) NOT NULL,
    "description" VARCHAR(255),
    "config" JSONB NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "integrations_pkey" PRIMARY KEY ("uuid")
);

-- CreateIndex
CREATE UNIQUE INDEX IF NOT EXISTS "integrations_name_key" ON "integrations"("name");

-- 2. Create shared_lists table
CREATE TABLE IF NOT EXISTS "shared_lists" (
    "name" VARCHAR(255) NOT NULL,
    "config" JSONB NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "shared_lists_pkey" PRIMARY KEY ("name")
);

-- 3. Alter nodes
ALTER TABLE "nodes" ADD COLUMN IF NOT EXISTS "integration_uuids" UUID[] DEFAULT ARRAY[]::UUID[];
ALTER TABLE "nodes" ADD COLUMN IF NOT EXISTS "ips" JSONB NOT NULL DEFAULT '[]';

-- 4. Alter hosts
ALTER TABLE "hosts" ADD COLUMN IF NOT EXISTS "mapper" JSONB NOT NULL DEFAULT '{}';

-- 5. Alter node_plugin
ALTER TABLE "node_plugin" ADD COLUMN IF NOT EXISTS "shared_lists" JSONB NOT NULL DEFAULT '[]';

-- 6. Drop legacy host columns
ALTER TABLE "hosts" DROP COLUMN IF EXISTS "singbox_custom_params";
ALTER TABLE "hosts" DROP COLUMN IF EXISTS "mihomo_custom_params";
ALTER TABLE "hosts" DROP COLUMN IF EXISTS "override_protocol_credential";
ALTER TABLE "hosts" DROP COLUMN IF EXISTS "protocol_credential";
ALTER TABLE "hosts" DROP COLUMN IF EXISTS "singbox_mux_params";
ALTER TABLE "hosts" DROP COLUMN IF EXISTS "clash_mux_params";
