"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInternalSquadByUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetInternalSquadByUuidCommand;
(function (GetInternalSquadByUuidCommand) {
    GetInternalSquadByUuidCommand.url = api_1.REST_API.INTERNAL_SQUADS.GET_BY_UUID;
    GetInternalSquadByUuidCommand.TSQ_url = GetInternalSquadByUuidCommand.url(':uuid');
    GetInternalSquadByUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.GET_BY_UUID(':uuid'), 'get', 'Get internal squad by uuid');
    GetInternalSquadByUuidCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetInternalSquadByUuidCommand.ResponseSchema = zod_1.z.object({
        response: models_1.InternalSquadSchema,
    });
})(GetInternalSquadByUuidCommand || (exports.GetInternalSquadByUuidCommand = GetInternalSquadByUuidCommand = {}));
