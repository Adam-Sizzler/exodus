-- Performance indexes backport from analysis vs upstream.
-- All created with IF NOT EXISTS — safe to run on any existing instance.

-- users.status: many hot queries filter by status = 'ACTIVE' / 'LIMITED' / 'EXPIRED'.
-- HASH is O(1) equality lookup, better than btree for this low-cardinality enum column.
CREATE INDEX IF NOT EXISTS users_status_idx
    ON public.users USING HASH (status);

-- users.tag: used in WHERE tag IS NOT NULL AND tag <> '' and ORDER BY tag ASC.
CREATE INDEX IF NOT EXISTS users_tag_idx
    ON public.users (tag);

-- nodes_user_usage_history: queried heavily by node (traffic aggregation per node).
-- Uses node_id (BigInt PK ref) instead of node_uuid — adapted for exodus schema.
CREATE INDEX IF NOT EXISTS nodes_user_usage_history_node_id_created_at_idx
    ON public.nodes_user_usage_history (node_id, created_at DESC);

-- nodes_user_usage_history: queried by user (per-user traffic history).
-- Uses user_id (BigInt PK ref to users.t_id) — adapted for exodus schema.
CREATE INDEX IF NOT EXISTS nodes_user_usage_history_user_id_updated_at_idx
    ON public.nodes_user_usage_history (user_id, updated_at DESC);

-- hwid_user_devices: PK is (hwid, user_uuid) which covers composite lookups.
-- This index covers single-column queries: WHERE user_uuid = ? and DELETE WHERE user_uuid = ?.
CREATE INDEX IF NOT EXISTS hwid_user_devices_user_uuid_idx
    ON public.hwid_user_devices (user_uuid);
