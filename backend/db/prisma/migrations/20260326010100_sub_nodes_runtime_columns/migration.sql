ALTER TABLE "sub_nodes"
    ADD COLUMN IF NOT EXISTS "last_status_change" TIMESTAMP(3),
    ADD COLUMN IF NOT EXISTS "last_status_message" TEXT,
    ADD COLUMN IF NOT EXISTS "singbox_version" TEXT,
    ADD COLUMN IF NOT EXISTS "node_version" TEXT,
    ADD COLUMN IF NOT EXISTS "singbox_uptime" TEXT NOT NULL DEFAULT '0',
    ADD COLUMN IF NOT EXISTS "cpu_count" INTEGER,
    ADD COLUMN IF NOT EXISTS "cpu_model" TEXT,
    ADD COLUMN IF NOT EXISTS "total_ram" TEXT;
