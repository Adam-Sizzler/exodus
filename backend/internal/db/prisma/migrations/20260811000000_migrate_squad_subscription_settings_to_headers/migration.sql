-- Migrate legacy per-squad subscription_settings fields (profile_title, support_link,
-- profile_update_interval, happ_announce, happ_routing, is_profile_webpage_url_enabled)
-- into response_headers_add, matching how the base subscription_settings table was already
-- migrated (see 20260802000000_api_schema). Upstream Exodus's per-squad override schema
-- only allows serve_json_at_base_subscription/is_show_custom_remarks/randomize_hosts to be
-- set via subscription_settings — everything else must go through responseHeadersAdd/Remove.
--
-- Existing response_headers_add keys always win (COALESCE order below) so an admin who
-- already configured the new-style header wins over a stale legacy field value.
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'external_squads' AND column_name = 'subscription_settings') THEN
        UPDATE "external_squads"
        SET "response_headers_add" =
            (
                (CASE WHEN subscription_settings ->> 'profile_title' IS NOT NULL AND subscription_settings ->> 'profile_title' <> ''
                     THEN jsonb_build_object('profile-title', 'exEncodeBase64:' || (subscription_settings ->> 'profile_title'))
                     ELSE '{}'::jsonb END)
                || (CASE WHEN subscription_settings ->> 'profile_update_interval' IS NOT NULL AND (subscription_settings ->> 'profile_update_interval')::int > 0
                     THEN jsonb_build_object('profile-update-interval', subscription_settings ->> 'profile_update_interval')
                     ELSE '{}'::jsonb END)
                || (CASE WHEN subscription_settings ->> 'support_link' IS NOT NULL AND subscription_settings ->> 'support_link' <> ''
                     THEN jsonb_build_object('support-url', subscription_settings ->> 'support_link')
                     ELSE '{}'::jsonb END)
                || (CASE WHEN subscription_settings ->> 'happ_announce' IS NOT NULL AND subscription_settings ->> 'happ_announce' <> ''
                     THEN jsonb_build_object('announce', 'exEncodeBase64:' || (subscription_settings ->> 'happ_announce'))
                     ELSE '{}'::jsonb END)
                || (CASE WHEN subscription_settings ->> 'happ_routing' IS NOT NULL AND subscription_settings ->> 'happ_routing' <> ''
                     THEN jsonb_build_object('routing', subscription_settings ->> 'happ_routing')
                     ELSE '{}'::jsonb END)
                || (CASE WHEN (subscription_settings ->> 'is_profile_webpage_url_enabled')::boolean IS TRUE
                     THEN jsonb_build_object('profile-web-page-url', '{{SUBSCRIPTION_URL}}')
                     ELSE '{}'::jsonb END)
            )
            || COALESCE("response_headers_add", '{}'::jsonb)
        WHERE subscription_settings IS NOT NULL AND jsonb_typeof(subscription_settings) = 'object';
    END IF;
END $$;
