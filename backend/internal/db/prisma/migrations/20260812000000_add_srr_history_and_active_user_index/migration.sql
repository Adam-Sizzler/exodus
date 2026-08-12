-- Add srr_response_type and srr_rule_name to user_subscription_request_history
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'user_subscription_request_history' AND column_name = 'srr_response_type') THEN
        ALTER TABLE "user_subscription_request_history" ADD COLUMN "srr_response_type" TEXT NOT NULL DEFAULT 'UNKNOWN';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'user_subscription_request_history' AND column_name = 'srr_rule_name') THEN
        ALTER TABLE "user_subscription_request_history" ADD COLUMN "srr_rule_name" TEXT;
    END IF;
END $$;

-- Add partial index for active users config lookup with all Exodus protocol credentials
CREATE INDEX IF NOT EXISTS "users_active_config_idx"
    ON "users" ("id")
    INCLUDE (
        "vless_uuid",
        "trojan_password",
        "ss_password",
        "naive_password",
        "shadowtls_password",
        "hysteria2_password",
        "anytls_password"
    )
    WHERE "status" = 'ACTIVE';

