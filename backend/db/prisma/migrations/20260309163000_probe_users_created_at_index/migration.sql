-- Probe migration after init snapshot.
-- This is a tiny forward-only change to validate the migration chain works.
CREATE INDEX "users_created_at_idx" ON "users"("created_at");
