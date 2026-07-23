"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateRemnawaveSettingsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var UpdateRemnawaveSettingsCommand;
(function (UpdateRemnawaveSettingsCommand) {
    UpdateRemnawaveSettingsCommand.url = api_1.REST_API.REMNAAWAVE_SETTINGS.UPDATE;
    UpdateRemnawaveSettingsCommand.TSQ_url = UpdateRemnawaveSettingsCommand.url;
    UpdateRemnawaveSettingsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.REMNAAWAVE_SETTINGS_ROUTES.UPDATE, 'patch', 'Update Remnawave settings', { scope: 'update', kind: 'write' });
    UpdateRemnawaveSettingsCommand.RequestSchema = zod_1.z.object({
        passkeySettings: models_1.PasskeySettingsSchema.optional(),
        oauth2Settings: models_1.Oauth2SettingsSchema.optional(),
        passwordSettings: models_1.PasswordAuthSettingsSchema.optional(),
        brandingSettings: models_1.BrandingSettingsSchema.optional(),
    });
    UpdateRemnawaveSettingsCommand.ResponseSchema = zod_1.z.object({
        response: models_1.RemnawaveSettingsSchema,
    });
})(UpdateRemnawaveSettingsCommand || (exports.UpdateRemnawaveSettingsCommand = UpdateRemnawaveSettingsCommand = {}));
