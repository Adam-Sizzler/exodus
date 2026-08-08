"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetInternalSquadCommand;
(function (GetInternalSquadCommand) {
    GetInternalSquadCommand.url = api_1.REST_API.INTERNAL_SQUADS.GET_BY_UUID;
    GetInternalSquadCommand.TSQ_url = GetInternalSquadCommand.url(':uuid');
    GetInternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.GET_BY_UUID(':uuid'), 'get', 'Get internal squad by uuid', { scope: 'get', kind: 'read' });
    GetInternalSquadCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    GetInternalSquadCommand.ResponseSchema = zod_1.z.object({
        response: models_1.InternalSquadSchema,
    });

})(GetInternalSquadCommand || (exports.GetInternalSquadCommand = GetInternalSquadCommand = {}));
