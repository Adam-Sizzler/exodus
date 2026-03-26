CREATE TABLE IF NOT EXISTS "sub_nodes" (
    "id" BIGSERIAL NOT NULL,
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "address" TEXT NOT NULL,
    "port" INTEGER,
    "api_schema" TEXT NOT NULL DEFAULT 'grpc',
    "api_path" TEXT NOT NULL DEFAULT '/',
    "is_connected" BOOLEAN NOT NULL DEFAULT false,
    "is_connecting" BOOLEAN NOT NULL DEFAULT false,
    "is_disabled" BOOLEAN NOT NULL DEFAULT false,
    "provider_uuid" UUID,
    "view_position" SERIAL NOT NULL,
    "tags" TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "sub_nodes_pkey" PRIMARY KEY ("uuid"),
    CONSTRAINT "sub_nodes_provider_uuid_fkey"
        FOREIGN KEY ("provider_uuid") REFERENCES "infra_providers"("uuid")
        ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS "sub_nodes_id_key" ON "sub_nodes"("id");
CREATE UNIQUE INDEX IF NOT EXISTS "sub_nodes_name_key" ON "sub_nodes"("name");
CREATE UNIQUE INDEX IF NOT EXISTS "sub_nodes_address_key" ON "sub_nodes"("address");
CREATE INDEX IF NOT EXISTS "sub_nodes_id_idx" ON "sub_nodes"("id");

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'subscription_connections'
    ) THEN
        INSERT INTO "sub_nodes" (
            "uuid",
            "name",
            "address",
            "port",
            "api_schema",
            "api_path",
            "is_connected",
            "is_connecting",
            "is_disabled",
            "provider_uuid",
            "view_position",
            "tags",
            "created_at",
            "updated_at"
        )
        SELECT
            "uuid",
            "name",
            "address",
            "port",
            COALESCE(NULLIF(TRIM("api_schema"), ''), 'grpc'),
            COALESCE(NULLIF(TRIM("api_path"), ''), '/'),
            COALESCE("is_connected", false),
            COALESCE("is_connecting", false),
            COALESCE("is_disabled", false),
            "provider_uuid",
            COALESCE("view_position", 1),
            COALESCE("tags", ARRAY[]::TEXT[]),
            COALESCE("created_at", now()),
            COALESCE("updated_at", now())
        FROM "subscription_connections"
        ON CONFLICT ("uuid") DO NOTHING;

        PERFORM setval(
            'sub_nodes_view_position_seq',
            GREATEST((SELECT COALESCE(MAX("view_position"), 0) FROM "sub_nodes") + 1, 1),
            false
        );
    END IF;
END$$;
