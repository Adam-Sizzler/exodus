ALTER TABLE public.nodes
    ADD COLUMN IF NOT EXISTS proxy_url text,
    ADD COLUMN IF NOT EXISTS node_consumption_multiplier bigint DEFAULT 1000000000 NOT NULL,
    ADD COLUMN IF NOT EXISTS note varchar(255);

UPDATE public.nodes
SET node_consumption_multiplier = COALESCE(node_consumption_multiplier, 1000000000)
WHERE node_consumption_multiplier IS NULL;

ALTER TABLE public.nodes
    ALTER COLUMN node_consumption_multiplier SET DEFAULT 1000000000,
    ALTER COLUMN node_consumption_multiplier SET NOT NULL;
