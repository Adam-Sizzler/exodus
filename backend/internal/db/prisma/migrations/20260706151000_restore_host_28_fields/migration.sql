ALTER TABLE "hosts" ADD COLUMN IF NOT EXISTS "final_mask" JSONB;
ALTER TABLE "hosts" ADD COLUMN IF NOT EXISTS "vless_route_id" INTEGER;
ALTER TABLE "hosts" ADD COLUMN IF NOT EXISTS "pinned_peer_cert_sha256" TEXT;
ALTER TABLE "hosts" ADD COLUMN IF NOT EXISTS "verify_peer_cert_by_name" TEXT;
