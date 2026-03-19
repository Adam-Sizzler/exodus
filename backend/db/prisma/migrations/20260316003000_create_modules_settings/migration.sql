CREATE TABLE IF NOT EXISTS "modules_settings" (
    "id" INTEGER NOT NULL DEFAULT 1,
    "haproxy_enabled" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "modules_settings_pkey" PRIMARY KEY ("id")
);
