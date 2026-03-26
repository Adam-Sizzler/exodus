ALTER TABLE "sub_nodes"
    ADD COLUMN IF NOT EXISTS "grpc_auth_token" TEXT;

UPDATE "sub_nodes"
SET "grpc_auth_token" = encode(gen_random_bytes(32), 'base64')
WHERE COALESCE(BTRIM("grpc_auth_token"), '') = '';

ALTER TABLE "sub_nodes"
    ALTER COLUMN "grpc_auth_token" SET DEFAULT encode(gen_random_bytes(32), 'base64');

ALTER TABLE "sub_nodes"
    ALTER COLUMN "grpc_auth_token" SET NOT NULL;
