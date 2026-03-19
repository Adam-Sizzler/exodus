"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInternalSquadsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetInternalSquadsCommand;
(function (GetInternalSquadsCommand) {
    GetInternalSquadsCommand.url = api_1.REST_API.INTERNAL_SQUADS.GET;
    GetInternalSquadsCommand.TSQ_url = GetInternalSquadsCommand.url;
    GetInternalSquadsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.GET, 'get', 'Get all internal squads');
    GetInternalSquadsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            internalSquads: zod_1.z.array(models_1.InternalSquadSchema),
        }),
    });
})(GetInternalSquadsCommand || (exports.GetInternalSquadsCommand = GetInternalSquadsCommand = {}));
