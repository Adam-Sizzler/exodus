"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateCerberusSettingsCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var UpdateCerberusSettingsCommand;
(function (UpdateCerberusSettingsCommand) {
    UpdateCerberusSettingsCommand.url = api_1.REST_API.REMNAAWAVE_SETTINGS.UPDATE;
    UpdateCerberusSettingsCommand.TSQ_url = UpdateCerberusSettingsCommand.url;
    UpdateCerberusSettingsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.REMNAAWAVE_SETTINGS_ROUTES.UPDATE, 'patch', 'Update Cerberus settings');
    UpdateCerberusSettingsCommand.RequestSchema = zod_1.z.object({
        passkeySettings: models_1.PasskeySettingsSchema.optional(),
        oauth2Settings: models_1.Oauth2SettingsSchema.optional(),
        tgAuthSettings: models_1.TgAuthSettingsSchema.optional(),
        passwordSettings: models_1.PasswordAuthSettingsSchema.optional(),
        brandingSettings: models_1.BrandingSettingsSchema.optional(),
    });
    UpdateCerberusSettingsCommand.ResponseSchema = zod_1.z.object({
        response: models_1.CerberusSettingsSchema,
    });
})(UpdateCerberusSettingsCommand || (exports.UpdateCerberusSettingsCommand = UpdateCerberusSettingsCommand = {}));
