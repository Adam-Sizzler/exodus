"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetExodusSettingsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetExodusSettingsCommand;
(function (GetExodusSettingsCommand) {
    GetExodusSettingsCommand.url = api_1.REST_API.REMNAAWAVE_SETTINGS.GET;
    GetExodusSettingsCommand.TSQ_url = GetExodusSettingsCommand.url;
    GetExodusSettingsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXODUS_SETTINGS_ROUTES.GET, 'get', 'Get Exodus settings', { scope: 'get', kind: 'read' });
    GetExodusSettingsCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExodusSettingsSchema,
    });

})(GetExodusSettingsCommand || (exports.GetExodusSettingsCommand = GetExodusSettingsCommand = {}));

exports.GetExodusSettingsCommand = exports.GetExodusSettingsCommand;
