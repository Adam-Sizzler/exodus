UPDATE public.users
SET trojan_password = BTRIM(trojan_password),
    ss_password = BTRIM(ss_password),
    naive_password = NULLIF(BTRIM(naive_password), ''),
    shadowtls_password = NULLIF(BTRIM(shadowtls_password), ''),
    hysteria2_password = NULLIF(BTRIM(hysteria2_password), ''),
    anytls_password = NULLIF(BTRIM(anytls_password), '');

UPDATE public.users
SET trojan_password = encode(gen_random_bytes(8), 'hex')
WHERE NULLIF(trojan_password, '') IS NULL
   OR trojan_password = vless_uuid::text;

UPDATE public.users
SET ss_password = encode(gen_random_bytes(8), 'hex')
WHERE NULLIF(ss_password, '') IS NULL
   OR ss_password IN (vless_uuid::text, trojan_password);

UPDATE public.users
SET naive_password = encode(gen_random_bytes(8), 'hex')
WHERE naive_password IS NULL
   OR naive_password IN (vless_uuid::text, trojan_password, ss_password);

UPDATE public.users
SET shadowtls_password = encode(gen_random_bytes(8), 'hex')
WHERE shadowtls_password IS NULL
   OR shadowtls_password IN (vless_uuid::text, trojan_password, ss_password, naive_password);

UPDATE public.users
SET hysteria2_password = encode(gen_random_bytes(8), 'hex')
WHERE hysteria2_password IS NULL
   OR hysteria2_password IN (vless_uuid::text, trojan_password, ss_password, naive_password, shadowtls_password);

UPDATE public.users
SET anytls_password = encode(gen_random_bytes(8), 'hex')
WHERE anytls_password IS NULL
   OR anytls_password IN (vless_uuid::text, trojan_password, ss_password, naive_password, shadowtls_password, hysteria2_password);

ALTER TABLE public.users
    ALTER COLUMN naive_password SET NOT NULL,
    ALTER COLUMN shadowtls_password SET NOT NULL,
    ALTER COLUMN hysteria2_password SET NOT NULL,
    ALTER COLUMN anytls_password SET NOT NULL;

ALTER TABLE public.users
    ADD CONSTRAINT users_protocol_credentials_unique CHECK (
        trojan_password <> vless_uuid::text
        AND ss_password <> vless_uuid::text
        AND naive_password <> vless_uuid::text
        AND shadowtls_password <> vless_uuid::text
        AND hysteria2_password <> vless_uuid::text
        AND anytls_password <> vless_uuid::text
        AND trojan_password <> ss_password
        AND trojan_password <> naive_password
        AND trojan_password <> shadowtls_password
        AND trojan_password <> hysteria2_password
        AND trojan_password <> anytls_password
        AND ss_password <> naive_password
        AND ss_password <> shadowtls_password
        AND ss_password <> hysteria2_password
        AND ss_password <> anytls_password
        AND naive_password <> shadowtls_password
        AND naive_password <> hysteria2_password
        AND naive_password <> anytls_password
        AND shadowtls_password <> hysteria2_password
        AND shadowtls_password <> anytls_password
        AND hysteria2_password <> anytls_password
    );
