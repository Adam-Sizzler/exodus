-- CreateTable
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

-- AlterTable
ALTER TABLE "nodes" ADD COLUMN IF NOT EXISTS "integration_uuids" UUID[] DEFAULT ARRAY[]::UUID[];

-- CreateTable
CREATE TABLE IF NOT EXISTS "shared_lists" (
    "name" VARCHAR(255) NOT NULL,
    "config" JSONB NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "shared_lists_pkey" PRIMARY KEY ("name")
);
