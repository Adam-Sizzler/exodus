ALTER TABLE "hosts"
    ADD COLUMN IF NOT EXISTS "singbox_mux_params" JSONB,
    ADD COLUMN IF NOT EXISTS "clash_mux_params" JSONB,
    ADD COLUMN IF NOT EXISTS "selector_nodes_first" BOOLEAN NOT NULL DEFAULT false;
