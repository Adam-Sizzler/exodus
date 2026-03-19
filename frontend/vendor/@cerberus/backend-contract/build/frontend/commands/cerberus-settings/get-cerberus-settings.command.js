"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetCerberusSettingsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
var GetCerberusSettingsCommand;
(function (GetCerberusSettingsCommand) {
    GetCerberusSettingsCommand.url = api_1.REST_API.REMNAAWAVE_SETTINGS.GET;
    GetCerberusSettingsCommand.TSQ_url = GetCerberusSettingsCommand.url;
    GetCerberusSettingsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.REMNAAWAVE_SETTINGS_ROUTES.GET, 'get', 'Get Cerberus settings');
    GetCerberusSettingsCommand.ResponseSchema = zod_1.z.object({
        response: models_1.CerberusSettingsSchema,
    });
})(GetCerberusSettingsCommand || (exports.GetCerberusSettingsCommand = GetCerberusSettingsCommand = {}));
