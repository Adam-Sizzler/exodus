"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetExternalSquadByUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetExternalSquadByUuidCommand;
(function (GetExternalSquadByUuidCommand) {
    GetExternalSquadByUuidCommand.url = api_1.REST_API.EXTERNAL_SQUADS.GET_BY_UUID;
    GetExternalSquadByUuidCommand.TSQ_url = GetExternalSquadByUuidCommand.url(':uuid');
    GetExternalSquadByUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXTERNAL_SQUADS_ROUTES.GET_BY_UUID(':uuid'), 'get', 'Get external squad by uuid');
    GetExternalSquadByUuidCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetExternalSquadByUuidCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExternalSquadSchema,
    });
})(GetExternalSquadByUuidCommand || (exports.GetExternalSquadByUuidCommand = GetExternalSquadByUuidCommand = {}));
