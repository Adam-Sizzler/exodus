-- Add SRS lists registry for sing-box ruleset management.
CREATE TABLE "srs_lists" (
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "type" TEXT NOT NULL DEFAULT 'remote',
    "tag" TEXT NOT NULL,
    "format" TEXT NOT NULL DEFAULT 'binary',
    "url" TEXT NOT NULL,
    "update_interval" TEXT NOT NULL DEFAULT '1d',
    "path" TEXT,
    "file_name" TEXT NOT NULL,
    "view_position" SERIAL NOT NULL,
    "is_available" BOOLEAN NOT NULL DEFAULT false,
    "last_checked_at" TIMESTAMP(3),
    "last_error" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT now(),
    CONSTRAINT "srs_lists_pkey" PRIMARY KEY ("uuid")
);

CREATE UNIQUE INDEX "srs_lists_tag_key" ON "srs_lists"("tag");
CREATE UNIQUE INDEX "srs_lists_url_key" ON "srs_lists"("url");
CREATE UNIQUE INDEX "srs_lists_file_name_key" ON "srs_lists"("file_name");
