ALTER TABLE "exodus_settings"
ADD COLUMN IF NOT EXISTS "modules_settings" JSONB;
