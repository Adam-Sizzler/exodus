-- clash_mux_params historically stores a Clash/Mihomo "smux" block, which is
-- native YAML, not JSON. It was kept as jsonb (like the sibling xray/singbox
-- mux columns, which genuinely are JSON) by wrapping the YAML text as a
-- JSON string scalar. That wrapping is the root cause of a class of bugs
-- across the frontend and generator code, so we drop it and store plain
-- YAML text instead.

-- 1. Add a temporary text column and unwrap existing jsonb-string values
--    into plain text. `#>> '{}'` returns the *unquoted, unescaped* text of a
--    jsonb scalar - unlike `::text`, which would keep the JSON quoting.
ALTER TABLE "hosts" ADD COLUMN "clash_mux_params_new" text;

UPDATE "hosts"
SET "clash_mux_params_new" = CASE
    WHEN "clash_mux_params" IS NULL THEN NULL
    -- Legacy rows that somehow ended up with a real JSON object instead of
    -- the wrapped-string convention: keep them as their raw JSON text so no
    -- data is silently dropped (parseMihomoMuxParams-equivalent code paths
    -- on read already tolerate a plain JSON object as input during rollout).
    WHEN jsonb_typeof("clash_mux_params") = 'string' THEN "clash_mux_params" #>> '{}'
    ELSE "clash_mux_params"::text
END;

-- 2. Swap columns.
ALTER TABLE "hosts" DROP COLUMN "clash_mux_params";
ALTER TABLE "hosts" RENAME COLUMN "clash_mux_params_new" TO "clash_mux_params";
