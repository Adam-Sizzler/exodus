ALTER TABLE public.hosts
  ADD COLUMN IF NOT EXISTS override_protocol_credential boolean DEFAULT false NOT NULL,
  ADD COLUMN IF NOT EXISTS protocol_credential text;

ALTER TABLE public.hosts
  DROP COLUMN IF EXISTS vless_route_id;
