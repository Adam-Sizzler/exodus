-- Drop legacy and custom host columns
ALTER TABLE hosts DROP COLUMN IF EXISTS singbox_custom_params;
ALTER TABLE hosts DROP COLUMN IF EXISTS mihomo_custom_params;
ALTER TABLE hosts DROP COLUMN IF EXISTS override_protocol_credential;
ALTER TABLE hosts DROP COLUMN IF EXISTS protocol_credential;
ALTER TABLE hosts DROP COLUMN IF EXISTS singbox_mux_params;
ALTER TABLE hosts DROP COLUMN IF EXISTS clash_mux_params;
