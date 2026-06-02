ALTER TABLE public.nodes
    ADD COLUMN IF NOT EXISTS grpc_auth_token text;

ALTER TABLE public.sub_nodes
    ADD COLUMN IF NOT EXISTS grpc_auth_token text;

UPDATE public.nodes
SET grpc_auth_token = COALESCE(
    (
        SELECT keygen.grpc_auth_token
        FROM public.keygen
        ORDER BY keygen.created_at ASC
        LIMIT 1
    ),
    encode(gen_random_bytes(32), 'hex'::text)
)
WHERE grpc_auth_token IS NULL
   OR btrim(grpc_auth_token) = '';

UPDATE public.sub_nodes
SET grpc_auth_token = COALESCE(
    (
        SELECT keygen.grpc_auth_token
        FROM public.keygen
        ORDER BY keygen.created_at ASC
        LIMIT 1
    ),
    encode(gen_random_bytes(32), 'hex'::text)
)
WHERE grpc_auth_token IS NULL
   OR btrim(grpc_auth_token) = '';

ALTER TABLE public.nodes
    ALTER COLUMN grpc_auth_token SET DEFAULT encode(gen_random_bytes(32), 'hex'::text),
    ALTER COLUMN grpc_auth_token SET NOT NULL;

ALTER TABLE public.sub_nodes
    ALTER COLUMN grpc_auth_token SET DEFAULT encode(gen_random_bytes(32), 'hex'::text),
    ALTER COLUMN grpc_auth_token SET NOT NULL;
