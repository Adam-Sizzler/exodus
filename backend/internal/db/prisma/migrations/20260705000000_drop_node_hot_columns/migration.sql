ALTER TABLE public.nodes
    DROP COLUMN IF EXISTS singbox_version,
    DROP COLUMN IF EXISTS node_version,
    DROP COLUMN IF EXISTS singbox_uptime,
    DROP COLUMN IF EXISTS users_online,
    DROP COLUMN IF EXISTS cpu_count,
    DROP COLUMN IF EXISTS cpu_model,
    DROP COLUMN IF EXISTS total_ram,
    DROP COLUMN IF EXISTS system_info,
    DROP COLUMN IF EXISTS system_stats;
