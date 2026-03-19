DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'nodes'
          AND column_name = 'xray_version'
    ) THEN
        ALTER TABLE "nodes" RENAME COLUMN "xray_version" TO "singbox_version";
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'nodes'
          AND column_name = 'xray_uptime'
    ) THEN
        ALTER TABLE "nodes" RENAME COLUMN "xray_uptime" TO "singbox_uptime";
    END IF;
END $$;
