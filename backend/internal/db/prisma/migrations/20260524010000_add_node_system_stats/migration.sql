ALTER TABLE public.nodes
    ADD COLUMN IF NOT EXISTS system_info jsonb,
    ADD COLUMN IF NOT EXISTS system_stats jsonb;
