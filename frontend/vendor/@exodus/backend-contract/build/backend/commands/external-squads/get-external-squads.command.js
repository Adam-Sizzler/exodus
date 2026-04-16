"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetExternalSquadsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetExternalSquadsCommand;
(function (GetExternalSquadsCommand) {
    GetExternalSquadsCommand.url = api_1.REST_API.EXTERNAL_SQUADS.GET;
    GetExternalSquadsCommand.TSQ_url = GetExternalSquadsCommand.url;
    GetExternalSquadsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXTERNAL_SQUADS_ROUTES.GET, 'get', 'Get all external squads');
    GetExternalSquadsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            externalSquads: zod_1.z.array(models_1.ExternalSquadSchema),
        }),
    });
})(GetExternalSquadsCommand || (exports.GetExternalSquadsCommand = GetExternalSquadsCommand = {}));
