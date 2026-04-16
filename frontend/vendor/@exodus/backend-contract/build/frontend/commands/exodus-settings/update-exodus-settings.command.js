"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateExodusSettingsCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var UpdateExodusSettingsCommand;
(function (UpdateExodusSettingsCommand) {
    UpdateExodusSettingsCommand.url = api_1.REST_API.EXODUS_SETTINGS.UPDATE;
    UpdateExodusSettingsCommand.TSQ_url = UpdateExodusSettingsCommand.url;
    UpdateExodusSettingsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXODUS_SETTINGS_ROUTES.UPDATE, 'patch', 'Update Exodus settings');
    UpdateExodusSettingsCommand.RequestSchema = zod_1.z.object({
        passkeySettings: models_1.PasskeySettingsSchema.optional(),
        oauth2Settings: models_1.Oauth2SettingsSchema.optional(),
        tgAuthSettings: models_1.TgAuthSettingsSchema.optional(),
        passwordSettings: models_1.PasswordAuthSettingsSchema.optional(),
        brandingSettings: models_1.BrandingSettingsSchema.optional(),
    });
    UpdateExodusSettingsCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExodusSettingsSchema,
    });
})(UpdateExodusSettingsCommand || (exports.UpdateExodusSettingsCommand = UpdateExodusSettingsCommand = {}));
