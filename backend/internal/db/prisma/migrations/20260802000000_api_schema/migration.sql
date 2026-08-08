-- 1. Rename t_id to id in users if t_id exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 't_id') THEN
        ALTER TABLE "users" RENAME COLUMN "t_id" TO "id";
    END IF;
END $$;

-- 2. Rename t_id to id in user_traffic if t_id exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'user_traffic' AND column_name = 't_id') THEN
        ALTER TABLE "user_traffic" RENAME COLUMN "t_id" TO "id";
    END IF;
END $$;

-- 3. Update external_squads for response_headers_add and response_headers_remove
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'external_squads' AND column_name = 'response_headers_add') THEN
        ALTER TABLE "external_squads" ADD COLUMN "response_headers_add" JSONB NOT NULL DEFAULT '{}';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'external_squads' AND column_name = 'response_headers_remove') THEN
        ALTER TABLE "external_squads" ADD COLUMN "response_headers_remove" TEXT[] DEFAULT ARRAY[]::TEXT[];
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'external_squads' AND column_name = 'response_headers') THEN
        UPDATE "external_squads" SET "response_headers_add" = "response_headers" WHERE "response_headers" IS NOT NULL AND jsonb_typeof("response_headers") = 'object';
        ALTER TABLE "external_squads" DROP COLUMN "response_headers";
    END IF;
END $$;

-- 4. Update subscription_settings for custom_response_headers migration
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'subscription_settings' AND column_name = 'profile_title') THEN
        UPDATE "subscription_settings"
        SET "custom_response_headers" =
            (
                jsonb_build_object(
                    'profile-title', 'exEncodeBase64:' || "profile_title",
                    'profile-update-interval', "profile_update_interval"::text,
                    'support-url', "support_link"
                )
                || (CASE WHEN "happ_announce" IS NOT NULL
                         THEN jsonb_build_object('announce', 'exEncodeBase64:' || "happ_announce")
                         ELSE '{}'::jsonb END)
                || (CASE WHEN "happ_routing" IS NOT NULL
                         THEN jsonb_build_object('routing', "happ_routing")
                         ELSE '{}'::jsonb END)
                || (CASE WHEN "is_profile_webpage_url_enabled"
                         THEN jsonb_build_object('profile-web-page-url', '{{SUBSCRIPTION_URL}}')
                         ELSE '{}'::jsonb END)
            )
            || COALESCE("custom_response_headers", '{}'::jsonb);

        ALTER TABLE "subscription_settings" DROP COLUMN IF EXISTS "happ_announce",
        DROP COLUMN IF EXISTS "happ_routing",
        DROP COLUMN IF EXISTS "is_profile_webpage_url_enabled",
        DROP COLUMN IF EXISTS "profile_title",
        DROP COLUMN IF EXISTS "profile_update_interval",
        DROP COLUMN IF EXISTS "support_link";
    END IF;
END $$;
