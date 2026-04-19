DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'cerberus_settings'
    ) AND NOT EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'exodus_settings'
    ) THEN
        ALTER TABLE "cerberus_settings" RENAME TO "exodus_settings";
    END IF;
END
$$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'cerberus_settings_pkey'
    ) THEN
        ALTER TABLE "exodus_settings"
            RENAME CONSTRAINT "cerberus_settings_pkey" TO "exodus_settings_pkey";
    END IF;
END
$$;
