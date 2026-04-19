-- This migration initializes the full schema from schema.prisma.

-- CreateTable
CREATE TABLE "cerberus_settings" (
    "id" INTEGER NOT NULL DEFAULT 1,
    "passkey_settings" JSONB,
    "oauth2_settings" JSONB,
    "tg_auth_settings" JSONB,
    "password_settings" JSONB,
    "branding_settings" JSONB,
    CONSTRAINT "cerberus_settings_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "users" (
    "t_id" BIGSERIAL NOT NULL,
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "short_uuid" TEXT NOT NULL,
    "username" TEXT NOT NULL,
    "status" VARCHAR(10) NOT NULL DEFAULT 'ACTIVE',
    "traffic_limit_bytes" BIGINT NOT NULL DEFAULT 0,
    "traffic_limit_strategy" TEXT NOT NULL DEFAULT 'NO_RESET',
    "expire_at" TIMESTAMP(3) NOT NULL,
    "sub_last_user_agent" TEXT,
    "sub_last_opened_at" TIMESTAMP(3),
    "last_traffic_reset_at" TIMESTAMP(3),
    "sub_revoked_at" TIMESTAMP(3),
    "trojan_password" TEXT NOT NULL,
    "vless_uuid" UUID NOT NULL,
    "ss_password" TEXT NOT NULL,
    "description" TEXT,
    "tag" TEXT,
    "telegram_id" BIGINT,
    "email" TEXT,
    "hwid_device_limit" INTEGER,
    "external_squad_uuid" UUID,
    "last_triggered_threshold" INTEGER NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "users_pkey" PRIMARY KEY ("t_id")
);

-- CreateTable
CREATE TABLE "user_traffic" (
    "t_id" BIGINT NOT NULL,
    "used_traffic_bytes" BIGINT NOT NULL DEFAULT 0,
    "lifetime_used_traffic_bytes" BIGINT NOT NULL DEFAULT 0,
    "online_at" TIMESTAMP(3),
    "last_connected_node_uuid" UUID,
    "first_connected_at" TIMESTAMP(3),
    CONSTRAINT "user_traffic_pkey" PRIMARY KEY ("t_id")
);

-- CreateTable
CREATE TABLE "api_tokens" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "token" TEXT NOT NULL,
    "token_name" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "api_tokens_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "admin" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "username" TEXT NOT NULL,
    "password_hash" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "session_ttl_minutes" INTEGER NOT NULL DEFAULT 60,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "admin_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "admin_sessions" (
    "session_token" TEXT NOT NULL,
    "admin_uuid" UUID NOT NULL,
    "expires_at" BIGINT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "admin_sessions_pkey" PRIMARY KEY ("session_token")
);

-- CreateTable
CREATE TABLE "passkeys" (
    "id" TEXT NOT NULL,
    "admin_uuid" UUID NOT NULL,
    "public_key" BYTEA NOT NULL,
    "counter" BIGINT NOT NULL,
    "device_type" TEXT NOT NULL,
    "backed_up" BOOLEAN NOT NULL,
    "transports" TEXT,
    "passkey_provider" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "passkeys_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "keygen" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "priv_key" TEXT NOT NULL,
    "pub_key" TEXT NOT NULL,
    "ca_cert" TEXT,
    "ca_key" TEXT,
    "client_cert" TEXT,
    "client_key" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "keygen_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "nodes" (
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
    "xray_version" TEXT,
    "node_version" TEXT,
    "xray_uptime" TEXT NOT NULL DEFAULT '0',
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
    CONSTRAINT "nodes_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "nodes_traffic_usage_history" (
    "id" BIGSERIAL NOT NULL,
    "node_uuid" UUID NOT NULL,
    "traffic_bytes" BIGINT NOT NULL,
    "reset_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "nodes_traffic_usage_history_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "nodes_user_usage_history" (
    "node_id" BIGINT NOT NULL,
    "user_id" BIGINT NOT NULL,
    "total_bytes" BIGINT NOT NULL,
    "created_at" DATE NOT NULL DEFAULT CURRENT_DATE,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "nodes_user_usage_history_pkey" PRIMARY KEY ("node_id", "created_at", "user_id")
);

-- CreateTable
CREATE TABLE "nodes_usage_history" (
    "node_uuid" UUID NOT NULL,
    "download_bytes" BIGINT NOT NULL,
    "upload_bytes" BIGINT NOT NULL,
    "total_bytes" BIGINT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT date_trunc('hour', now()),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "nodes_usage_history_pkey" PRIMARY KEY ("node_uuid", "created_at")
);

-- CreateTable
CREATE TABLE "hosts" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "view_position" SERIAL NOT NULL,
    "remark" VARCHAR(50) NOT NULL,
    "address" TEXT NOT NULL,
    "port" INTEGER NOT NULL,
    "path" TEXT,
    "sni" TEXT,
    "host" TEXT,
    "alpn" TEXT,
    "fingerprint" TEXT,
    "security_layer" TEXT NOT NULL DEFAULT 'DEFAULT',
    "xhttp_extra_params" JSONB,
    "mux_params" JSONB,
    "sockopt_params" JSONB,
    "is_disabled" BOOLEAN NOT NULL DEFAULT false,
    "server_description" VARCHAR(30),
    "vless_route_id" INTEGER,
    "allow_insecure" BOOLEAN NOT NULL DEFAULT false,
    "shuffle_host" BOOLEAN NOT NULL DEFAULT false,
    "mihomo_x25519" BOOLEAN NOT NULL DEFAULT false,
    "xray_json_template_uuid" UUID,
    "keep_sni_blank" BOOLEAN NOT NULL DEFAULT false,
    "exclude_from_subscription_types" TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    "tag" TEXT,
    "is_hidden" BOOLEAN NOT NULL DEFAULT false,
    "override_sni_from_address" BOOLEAN NOT NULL DEFAULT false,
    "config_profile_uuid" UUID,
    "config_profile_inbound_uuid" UUID,
    CONSTRAINT "hosts_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "internal_squad_host_exclusions" (
    "host_uuid" UUID NOT NULL,
    "squad_uuid" UUID NOT NULL,
    CONSTRAINT "internal_squad_host_exclusions_pkey" PRIMARY KEY ("host_uuid", "squad_uuid")
);

-- CreateTable
CREATE TABLE "subscription_templates" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "view_position" SERIAL NOT NULL,
    "name" VARCHAR(255) NOT NULL DEFAULT 'Default',
    "template_type" TEXT NOT NULL,
    "template_yaml" TEXT,
    "template_json" JSONB,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "subscription_templates_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "subscription_settings" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "profile_title" TEXT NOT NULL,
    "support_link" TEXT NOT NULL,
    "profile_update_interval" INTEGER NOT NULL,
    "address" TEXT,
    "port" INTEGER,
    "api_schema" TEXT,
    "api_path" TEXT,
    "is_profile_webpage_url_enabled" BOOLEAN NOT NULL DEFAULT true,
    "serve_json_at_base_subscription" BOOLEAN NOT NULL DEFAULT false,
    "happ_announce" TEXT,
    "happ_routing" TEXT,
    "is_show_custom_remarks" BOOLEAN NOT NULL DEFAULT true,
    "custom_remarks" JSONB NOT NULL,
    "custom_response_headers" JSONB,
    "randomize_hosts" BOOLEAN NOT NULL DEFAULT false,
    "response_rules" JSONB,
    "hwid_settings" JSONB,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "subscription_settings_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "hwid_user_devices" (
    "hwid" TEXT NOT NULL,
    "user_uuid" UUID NOT NULL,
    "platform" TEXT,
    "os_version" TEXT,
    "device_model" TEXT,
    "user_agent" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "hwid_user_devices_pkey" PRIMARY KEY ("hwid", "user_uuid")
);

-- CreateTable
CREATE TABLE "internal_squads" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "view_position" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "internal_squads_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "internal_squad_members" (
    "internal_squad_uuid" UUID NOT NULL,
    "user_id" BIGINT NOT NULL,
    CONSTRAINT "internal_squad_members_pkey" PRIMARY KEY ("internal_squad_uuid", "user_id")
);

-- CreateTable
CREATE TABLE "internal_squad_inbounds" (
    "internal_squad_uuid" UUID NOT NULL,
    "inbound_uuid" UUID NOT NULL,
    CONSTRAINT "internal_squad_inbounds_pkey" PRIMARY KEY ("internal_squad_uuid", "inbound_uuid")
);

-- CreateTable
CREATE TABLE "config_profiles" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "view_position" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "config" JSONB NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "config_profiles_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "config_profile_inbounds" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "profile_uuid" UUID NOT NULL,
    "tag" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    "network" TEXT,
    "security" TEXT,
    "port" INTEGER,
    "raw_inbound" JSONB,
    CONSTRAINT "config_profile_inbounds_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "config_profile_inbounds_to_nodes" (
    "config_profile_inbound_uuid" UUID NOT NULL,
    "node_uuid" UUID NOT NULL,
    CONSTRAINT "config_profile_inbounds_to_nodes_pkey" PRIMARY KEY ("config_profile_inbound_uuid", "node_uuid")
);

-- CreateTable
CREATE TABLE "infra_providers" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "favicon_link" TEXT,
    "login_url" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "infra_providers_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "infra_billing_nodes" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "node_uuid" UUID NOT NULL,
    "provider_uuid" UUID NOT NULL,
    "next_billing_at" TIMESTAMP(3) NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "infra_billing_nodes_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "infra_billing_history" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "provider_uuid" UUID NOT NULL,
    "amount" DOUBLE PRECISION NOT NULL,
    "billed_at" TIMESTAMP(3) NOT NULL,
    CONSTRAINT "infra_billing_history_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "user_subscription_request_history" (
    "id" BIGSERIAL NOT NULL,
    "user_uuid" UUID NOT NULL,
    "request_ip" TEXT,
    "user_agent" TEXT,
    "request_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "user_subscription_request_history_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "hosts_to_nodes" (
    "host_uuid" UUID NOT NULL,
    "node_uuid" UUID NOT NULL,
    CONSTRAINT "hosts_to_nodes_pkey" PRIMARY KEY ("host_uuid", "node_uuid")
);

-- CreateTable
CREATE TABLE "config_profile_snippets" (
    "name" VARCHAR(255) NOT NULL,
    "snippet" JSONB NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "config_profile_snippets_pkey" PRIMARY KEY ("name")
);

-- CreateTable
CREATE TABLE "external_squads" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "view_position" SERIAL NOT NULL,
    "name" VARCHAR(30) NOT NULL,
    "subscription_settings" JSONB,
    "host_overrides" JSONB,
    "response_headers" JSONB,
    "hwid_settings" JSONB,
    "custom_remarks" JSONB,
    "subpage_config_uuid" UUID,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "external_squads_pkey" PRIMARY KEY ("uuid")
);

-- CreateTable
CREATE TABLE "external_squads_templates" (
    "external_squad_uuid" UUID NOT NULL,
    "template_uuid" UUID NOT NULL,
    "template_type" TEXT NOT NULL,
    CONSTRAINT "external_squads_templates_pkey" PRIMARY KEY ("external_squad_uuid", "template_type")
);

-- CreateTable
CREATE TABLE "subscription_page_config" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "view_position" SERIAL NOT NULL,
    "name" VARCHAR(30) NOT NULL,
    "config" JSONB NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "subscription_page_config_pkey" PRIMARY KEY ("uuid")
);

-- CreateIndex
CREATE UNIQUE INDEX "users_uuid_key" ON "users"("uuid");

-- CreateIndex
CREATE UNIQUE INDEX "users_short_uuid_key" ON "users"("short_uuid");

-- CreateIndex
CREATE UNIQUE INDEX "users_username_key" ON "users"("username");

-- CreateIndex
CREATE UNIQUE INDEX "api_tokens_token_key" ON "api_tokens"("token");

-- CreateIndex
CREATE UNIQUE INDEX "admin_username_key" ON "admin"("username");

-- CreateIndex
CREATE INDEX "admin_sessions_admin_uuid_idx" ON "admin_sessions"("admin_uuid");

-- CreateIndex
CREATE INDEX "admin_sessions_expires_at_idx" ON "admin_sessions"("expires_at");

-- CreateIndex
CREATE INDEX "passkeys_id_idx" ON "passkeys"("id");

-- CreateIndex
CREATE INDEX "passkeys_admin_uuid_idx" ON "passkeys"("admin_uuid");

-- CreateIndex
CREATE UNIQUE INDEX "nodes_id_key" ON "nodes"("id");

-- CreateIndex
CREATE UNIQUE INDEX "nodes_name_key" ON "nodes"("name");

-- CreateIndex
CREATE UNIQUE INDEX "nodes_address_key" ON "nodes"("address");

-- CreateIndex
CREATE INDEX "nodes_id_idx" ON "nodes"("id");

-- CreateIndex
CREATE INDEX "nodes_usage_history_node_uuid_created_at_idx" ON "nodes_usage_history"("node_uuid", "created_at" DESC);

-- CreateIndex
CREATE UNIQUE INDEX "subscription_templates_template_type_name_key" ON "subscription_templates"("template_type", "name");

-- CreateIndex
CREATE UNIQUE INDEX "internal_squads_name_key" ON "internal_squads"("name");

-- CreateIndex
CREATE INDEX "internal_squad_members_internal_squad_uuid_idx" ON "internal_squad_members"("internal_squad_uuid");

-- CreateIndex
CREATE INDEX "internal_squad_members_user_id_idx" ON "internal_squad_members"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "config_profiles_name_key" ON "config_profiles"("name");

-- CreateIndex
CREATE UNIQUE INDEX "config_profile_inbounds_tag_key" ON "config_profile_inbounds"("tag");

-- CreateIndex
CREATE INDEX "config_profile_inbounds_profile_uuid_uuid_idx" ON "config_profile_inbounds"("profile_uuid", "uuid");

-- CreateIndex
CREATE UNIQUE INDEX "infra_providers_name_key" ON "infra_providers"("name");

-- CreateIndex
CREATE UNIQUE INDEX "infra_billing_nodes_node_uuid_provider_uuid_key" ON "infra_billing_nodes"("node_uuid", "provider_uuid");

-- CreateIndex
CREATE INDEX "infra_billing_nodes_next_billing_at_idx" ON "infra_billing_nodes"("next_billing_at");

-- CreateIndex
CREATE INDEX "user_subscription_request_history_user_uuid_idx" ON "user_subscription_request_history"("user_uuid");

-- CreateIndex
CREATE INDEX "user_subscription_request_history_request_at_idx" ON "user_subscription_request_history"("request_at" ASC);

-- CreateIndex
CREATE UNIQUE INDEX "external_squads_name_key" ON "external_squads"("name");

-- CreateIndex
CREATE UNIQUE INDEX "subscription_page_config_name_key" ON "subscription_page_config"("name");

-- AddForeignKey
ALTER TABLE "users" ADD CONSTRAINT "users_external_squad_uuid_fkey" FOREIGN KEY ("external_squad_uuid") REFERENCES "external_squads"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "user_traffic" ADD CONSTRAINT "user_traffic_t_id_fkey" FOREIGN KEY ("t_id") REFERENCES "users"("t_id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "user_traffic" ADD CONSTRAINT "user_traffic_last_connected_node_uuid_fkey" FOREIGN KEY ("last_connected_node_uuid") REFERENCES "nodes"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "admin_sessions" ADD CONSTRAINT "admin_sessions_admin_uuid_fkey" FOREIGN KEY ("admin_uuid") REFERENCES "admin"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "passkeys" ADD CONSTRAINT "passkeys_admin_uuid_fkey" FOREIGN KEY ("admin_uuid") REFERENCES "admin"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "nodes" ADD CONSTRAINT "nodes_active_config_profile_uuid_fkey" FOREIGN KEY ("active_config_profile_uuid") REFERENCES "config_profiles"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "nodes" ADD CONSTRAINT "nodes_provider_uuid_fkey" FOREIGN KEY ("provider_uuid") REFERENCES "infra_providers"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "nodes_traffic_usage_history" ADD CONSTRAINT "nodes_traffic_usage_history_node_uuid_fkey" FOREIGN KEY ("node_uuid") REFERENCES "nodes"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "nodes_user_usage_history" ADD CONSTRAINT "nodes_user_usage_history_node_id_fkey" FOREIGN KEY ("node_id") REFERENCES "nodes"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "nodes_user_usage_history" ADD CONSTRAINT "nodes_user_usage_history_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("t_id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "nodes_usage_history" ADD CONSTRAINT "nodes_usage_history_node_uuid_fkey" FOREIGN KEY ("node_uuid") REFERENCES "nodes"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "hosts" ADD CONSTRAINT "hosts_config_profile_inbound_uuid_fkey" FOREIGN KEY ("config_profile_inbound_uuid") REFERENCES "config_profile_inbounds"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "hosts" ADD CONSTRAINT "hosts_config_profile_uuid_fkey" FOREIGN KEY ("config_profile_uuid") REFERENCES "config_profiles"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "hosts" ADD CONSTRAINT "hosts_xray_json_template_uuid_fkey" FOREIGN KEY ("xray_json_template_uuid") REFERENCES "subscription_templates"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "internal_squad_host_exclusions" ADD CONSTRAINT "internal_squad_host_exclusions_host_uuid_fkey" FOREIGN KEY ("host_uuid") REFERENCES "hosts"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "internal_squad_host_exclusions" ADD CONSTRAINT "internal_squad_host_exclusions_squad_uuid_fkey" FOREIGN KEY ("squad_uuid") REFERENCES "internal_squads"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "hwid_user_devices" ADD CONSTRAINT "hwid_user_devices_user_uuid_fkey" FOREIGN KEY ("user_uuid") REFERENCES "users"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "internal_squad_members" ADD CONSTRAINT "internal_squad_members_internal_squad_uuid_fkey" FOREIGN KEY ("internal_squad_uuid") REFERENCES "internal_squads"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "internal_squad_members" ADD CONSTRAINT "internal_squad_members_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("t_id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "internal_squad_inbounds" ADD CONSTRAINT "internal_squad_inbounds_internal_squad_uuid_fkey" FOREIGN KEY ("internal_squad_uuid") REFERENCES "internal_squads"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "internal_squad_inbounds" ADD CONSTRAINT "internal_squad_inbounds_inbound_uuid_fkey" FOREIGN KEY ("inbound_uuid") REFERENCES "config_profile_inbounds"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "config_profile_inbounds" ADD CONSTRAINT "config_profile_inbounds_profile_uuid_fkey" FOREIGN KEY ("profile_uuid") REFERENCES "config_profiles"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "config_profile_inbounds_to_nodes" ADD CONSTRAINT "config_profile_inbounds_to_nodes_config_profile_inbound_uuid_fkey" FOREIGN KEY ("config_profile_inbound_uuid") REFERENCES "config_profile_inbounds"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "config_profile_inbounds_to_nodes" ADD CONSTRAINT "config_profile_inbounds_to_nodes_node_uuid_fkey" FOREIGN KEY ("node_uuid") REFERENCES "nodes"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "infra_billing_nodes" ADD CONSTRAINT "infra_billing_nodes_node_uuid_fkey" FOREIGN KEY ("node_uuid") REFERENCES "nodes"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "infra_billing_nodes" ADD CONSTRAINT "infra_billing_nodes_provider_uuid_fkey" FOREIGN KEY ("provider_uuid") REFERENCES "infra_providers"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "infra_billing_history" ADD CONSTRAINT "infra_billing_history_provider_uuid_fkey" FOREIGN KEY ("provider_uuid") REFERENCES "infra_providers"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "user_subscription_request_history" ADD CONSTRAINT "user_subscription_request_history_user_uuid_fkey" FOREIGN KEY ("user_uuid") REFERENCES "users"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "hosts_to_nodes" ADD CONSTRAINT "hosts_to_nodes_host_uuid_fkey" FOREIGN KEY ("host_uuid") REFERENCES "hosts"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "hosts_to_nodes" ADD CONSTRAINT "hosts_to_nodes_node_uuid_fkey" FOREIGN KEY ("node_uuid") REFERENCES "nodes"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "external_squads" ADD CONSTRAINT "external_squads_subpage_config_uuid_fkey" FOREIGN KEY ("subpage_config_uuid") REFERENCES "subscription_page_config"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "external_squads_templates" ADD CONSTRAINT "external_squads_templates_external_squad_uuid_fkey" FOREIGN KEY ("external_squad_uuid") REFERENCES "external_squads"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "external_squads_templates" ADD CONSTRAINT "external_squads_templates_template_uuid_fkey" FOREIGN KEY ("template_uuid") REFERENCES "subscription_templates"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;
