DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'config_profiles'
          AND column_name = 'config'
          AND udt_name = 'jsonb'
    ) THEN
        ALTER TABLE "config_profiles"
            ALTER COLUMN "config" TYPE JSON USING "config"::json;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'subscription_templates'
          AND column_name = 'template_json'
          AND udt_name = 'jsonb'
    ) THEN
        ALTER TABLE "subscription_templates"
            ALTER COLUMN "template_json" TYPE JSON USING "template_json"::json;
    END IF;
END
$$;
