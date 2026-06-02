ALTER TABLE public.hosts
    ADD COLUMN IF NOT EXISTS xhttp_extra_params jsonb;
