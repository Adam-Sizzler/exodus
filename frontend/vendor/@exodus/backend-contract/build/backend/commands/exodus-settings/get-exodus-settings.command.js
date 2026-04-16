"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetExodusSettingsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
var GetExodusSettingsCommand;
(function (GetExodusSettingsCommand) {
    GetExodusSettingsCommand.url = api_1.REST_API.EXODUS_SETTINGS.GET;
    GetExodusSettingsCommand.TSQ_url = GetExodusSettingsCommand.url;
    GetExodusSettingsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXODUS_SETTINGS_ROUTES.GET, 'get', 'Get Exodus settings');
    GetExodusSettingsCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExodusSettingsSchema,
    });
})(GetExodusSettingsCommand || (exports.GetExodusSettingsCommand = GetExodusSettingsCommand = {}));
