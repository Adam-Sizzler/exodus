ALTER TABLE "hosts" ADD COLUMN IF NOT EXISTS "tags" TEXT[] DEFAULT ARRAY[]::TEXT[];

UPDATE "hosts"
SET "tags" = ARRAY["tag"]
WHERE "tag" IS NOT NULL
  AND "tag" <> ''
  AND ("tags" IS NULL OR cardinality("tags") = 0);

ALTER TABLE "hosts" ALTER COLUMN "tags" SET DEFAULT ARRAY[]::TEXT[];
ALTER TABLE "hosts" ALTER COLUMN "tags" SET NOT NULL;
ALTER TABLE "hosts" DROP COLUMN IF EXISTS "tag";
