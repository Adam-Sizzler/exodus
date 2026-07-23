"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetRemnawaveSettingsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetRemnawaveSettingsCommand;
(function (GetRemnawaveSettingsCommand) {
    GetRemnawaveSettingsCommand.url = api_1.REST_API.REMNAAWAVE_SETTINGS.GET;
    GetRemnawaveSettingsCommand.TSQ_url = GetRemnawaveSettingsCommand.url;
    GetRemnawaveSettingsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.REMNAAWAVE_SETTINGS_ROUTES.GET, 'get', 'Get Remnawave settings', { scope: 'get', kind: 'read' });
    GetRemnawaveSettingsCommand.ResponseSchema = zod_1.z.object({
        response: models_1.RemnawaveSettingsSchema,
    });
})(GetRemnawaveSettingsCommand || (exports.GetRemnawaveSettingsCommand = GetRemnawaveSettingsCommand = {}));
