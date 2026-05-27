ALTER TABLE users
    ADD COLUMN IF NOT EXISTS naive_password TEXT,
    ADD COLUMN IF NOT EXISTS shadowtls_password TEXT,
    ADD COLUMN IF NOT EXISTS hysteria2_password TEXT,
    ADD COLUMN IF NOT EXISTS anytls_password TEXT;
