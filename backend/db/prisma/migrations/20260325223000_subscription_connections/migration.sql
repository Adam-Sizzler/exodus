CREATE TABLE IF NOT EXISTS "subscription_connections" (
    "id" BIGSERIAL NOT NULL,
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "address" TEXT NOT NULL,
    "port" INTEGER,
    "api_schema" TEXT NOT NULL DEFAULT 'grpc',
    "api_path" TEXT NOT NULL DEFAULT '',
    "active_config_profile_uuid" UUID,
    "is_connected" BOOLEAN NOT NULL DEFAULT false,
    "is_connecting" BOOLEAN NOT NULL DEFAULT false,
    "is_disabled" BOOLEAN NOT NULL DEFAULT false,
    "last_status_change" TIMESTAMP(3),
    "last_status_message" TEXT,
    "singbox_version" TEXT,
    "node_version" TEXT,
    "singbox_uptime" TEXT NOT NULL DEFAULT '0',
    "users_online" INTEGER DEFAULT 0,
    "consumption_multiplier" BIGINT NOT NULL DEFAULT 1000000000,
    "is_traffic_tracking_active" BOOLEAN NOT NULL DEFAULT false,
    "traffic_reset_day" INTEGER DEFAULT 1,
    "traffic_limit_bytes" BIGINT DEFAULT 0,
    "traffic_used_bytes" BIGINT DEFAULT 0,
    "notify_percent" INTEGER DEFAULT 0,
    "provider_uuid" UUID,
    "view_position" SERIAL NOT NULL,
    "country_code" TEXT NOT NULL DEFAULT 'XX',
    "tags" TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    "cpu_count" INTEGER,
    "cpu_model" TEXT,
    "total_ram" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "subscription_connections_pkey" PRIMARY KEY ("uuid"),
    CONSTRAINT "subscription_connections_active_config_profile_uuid_fkey"
        FOREIGN KEY ("active_config_profile_uuid") REFERENCES "config_profiles"("uuid")
        ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT "subscription_connections_provider_uuid_fkey"
        FOREIGN KEY ("provider_uuid") REFERENCES "infra_providers"("uuid")
        ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS "subscription_connections_traffic_usage_history" (
    "id" BIGSERIAL NOT NULL,
    "node_uuid" UUID NOT NULL,
    "traffic_bytes" BIGINT NOT NULL,
    "reset_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "subscription_connections_traffic_usage_history_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "subscription_connections_traffic_usage_history_node_uuid_fkey"
        FOREIGN KEY ("node_uuid") REFERENCES "subscription_connections"("uuid")
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS "config_profile_inbounds_to_subscription_connections" (
    "config_profile_inbound_uuid" UUID NOT NULL,
    "node_uuid" UUID NOT NULL,
    CONSTRAINT "config_profile_inbounds_to_subscription_connections_pkey"
        PRIMARY KEY ("config_profile_inbound_uuid", "node_uuid"),
    CONSTRAINT "config_profile_inbounds_to_subscription_connections_config_profile_inbound_uuid_fkey"
        FOREIGN KEY ("config_profile_inbound_uuid") REFERENCES "config_profile_inbounds"("uuid")
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT "config_profile_inbounds_to_subscription_connections_node_uuid_fkey"
        FOREIGN KEY ("node_uuid") REFERENCES "subscription_connections"("uuid")
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS "subscription_connections_id_key"
    ON "subscription_connections"("id");
CREATE UNIQUE INDEX IF NOT EXISTS "subscription_connections_name_key"
    ON "subscription_connections"("name");
CREATE UNIQUE INDEX IF NOT EXISTS "subscription_connections_address_key"
    ON "subscription_connections"("address");
CREATE INDEX IF NOT EXISTS "subscription_connections_id_idx"
    ON "subscription_connections"("id");
