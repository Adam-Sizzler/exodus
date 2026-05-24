ALTER TABLE public.nodes
    ADD COLUMN IF NOT EXISTS active_plugin_uuid uuid;
