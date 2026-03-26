ALTER TABLE IF EXISTS "keygen"
    ADD COLUMN IF NOT EXISTS "grpc_auth_token" TEXT;

UPDATE "keygen"
SET "grpc_auth_token" = encode(gen_random_bytes(32), 'hex')
WHERE COALESCE(BTRIM("grpc_auth_token"), '') = '';

ALTER TABLE IF EXISTS "keygen"
    ALTER COLUMN "grpc_auth_token" SET DEFAULT encode(gen_random_bytes(32), 'hex');

ALTER TABLE IF EXISTS "keygen"
    ALTER COLUMN "grpc_auth_token" SET NOT NULL;

UPDATE "sub_nodes"
SET "api_schema" = CASE
    WHEN LOWER(BTRIM("api_schema")) IN ('tls', 'https', 'grpcs') THEN 'tls'
    ELSE 'mtls'
END;

ALTER TABLE IF EXISTS "sub_nodes"
    ALTER COLUMN "api_schema" SET DEFAULT 'mtls';

ALTER TABLE IF EXISTS "sub_nodes"
    DROP COLUMN IF EXISTS "grpc_auth_token",
    DROP COLUMN IF EXISTS "singbox_version",
    DROP COLUMN IF EXISTS "node_version",
    DROP COLUMN IF EXISTS "singbox_uptime",
    DROP COLUMN IF EXISTS "cpu_count",
    DROP COLUMN IF EXISTS "cpu_model",
    DROP COLUMN IF EXISTS "total_ram";
