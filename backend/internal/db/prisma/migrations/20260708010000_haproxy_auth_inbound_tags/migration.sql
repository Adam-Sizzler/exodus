ALTER TABLE public.node_plugins
    ALTER COLUMN plugin_config SET DEFAULT '{"ingressFilter":{"enabled":false,"blockedIps":[]},"egressFilter":{"enabled":false,"blockedIps":[],"blockedPorts":[]},"haproxyAuth":{"inboundTags":[]},"sharedLists":[]}'::jsonb;

WITH plugin_tag_rows AS (
    SELECT DISTINCT
        np.uuid AS plugin_uuid,
        btrim(cpi.tag) AS tag
    FROM public.node_plugins np
    JOIN public.nodes n ON n.active_plugin_uuid = np.uuid
    JOIN public.config_profile_inbounds_to_nodes cpitn ON cpitn.node_uuid = n.uuid
    JOIN public.config_profile_inbounds cpi ON cpi.uuid = cpitn.config_profile_inbound_uuid
    WHERE lower(cpi.type) IN ('vless', 'trojan', 'naive', 'anytls')
      AND btrim(COALESCE(cpi.tag, '')) <> ''
), plugin_tags AS (
    SELECT
        plugin_uuid,
        jsonb_agg(tag ORDER BY tag) AS inbound_tags
    FROM plugin_tag_rows
    GROUP BY plugin_uuid
), all_plugin_tags AS (
    SELECT
        np.uuid AS plugin_uuid,
        COALESCE(pt.inbound_tags, '[]'::jsonb) AS inbound_tags
    FROM public.node_plugins np
    LEFT JOIN plugin_tags pt ON pt.plugin_uuid = np.uuid
)
UPDATE public.node_plugins np
SET plugin_config = jsonb_set(
    COALESCE(np.plugin_config, '{}'::jsonb) - 'haproxyAuth',
    '{haproxyAuth}',
    jsonb_build_object(
        'inboundTags',
        CASE
            WHEN jsonb_typeof(np.plugin_config -> 'haproxyAuth' -> 'inboundTags') = 'array'
                THEN np.plugin_config -> 'haproxyAuth' -> 'inboundTags'
            WHEN jsonb_typeof(np.plugin_config -> 'haproxyAuth' -> 'enabled') = 'boolean'
                 AND (np.plugin_config -> 'haproxyAuth' ->> 'enabled')::boolean
                THEN apt.inbound_tags
            ELSE '[]'::jsonb
        END
    ),
    true
)
FROM all_plugin_tags apt
WHERE np.uuid = apt.plugin_uuid;
