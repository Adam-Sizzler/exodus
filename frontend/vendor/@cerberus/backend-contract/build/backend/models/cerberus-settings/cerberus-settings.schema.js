"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CerberusSettingsSchema = void 0;
const zod_1 = require("zod");
const password_auth_settings_schema_1 = require("./password-auth-settings.schema");
const branding_settings_schema_1 = require("./branding-settings.schema");
const passkey_settings_schema_1 = require("./passkey-settings.schema");
const tg_auth_settings_schema_1 = require("./tg-auth-settings.schema");
const oauth2_settings_schema_1 = require("./oauth2-settings.schema");
exports.CerberusSettingsSchema = zod_1.z.object({
    passkeySettings: zod_1.z.nullable(passkey_settings_schema_1.PasskeySettingsSchema),
    oauth2Settings: zod_1.z.nullable(oauth2_settings_schema_1.Oauth2SettingsSchema),
    tgAuthSettings: zod_1.z.nullable(tg_auth_settings_schema_1.TgAuthSettingsSchema),
    passwordSettings: zod_1.z.nullable(password_auth_settings_schema_1.PasswordAuthSettingsSchema),
    brandingSettings: zod_1.z.nullable(branding_settings_schema_1.BrandingSettingsSchema),
});
