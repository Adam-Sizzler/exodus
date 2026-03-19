"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInternalSquadAccessibleNodesCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetInternalSquadAccessibleNodesCommand;
(function (GetInternalSquadAccessibleNodesCommand) {
    GetInternalSquadAccessibleNodesCommand.url = api_1.REST_API.INTERNAL_SQUADS.ACCESSIBLE_NODES;
    GetInternalSquadAccessibleNodesCommand.TSQ_url = GetInternalSquadAccessibleNodesCommand.url(':uuid');
    GetInternalSquadAccessibleNodesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.ACCESSIBLE_NODES(':uuid'), 'get', 'Get internal squad accessible nodes');
    GetInternalSquadAccessibleNodesCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetInternalSquadAccessibleNodesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            squadUuid: zod_1.z.string().uuid(),
            accessibleNodes: zod_1.z.array(zod_1.z.object({
                uuid: zod_1.z.string().uuid(),
                nodeName: zod_1.z.string(),
                countryCode: zod_1.z.string(),
                configProfileUuid: zod_1.z.string().uuid(),
                configProfileName: zod_1.z.string(),
                activeInbounds: zod_1.z.array(zod_1.z.string()),
            })),
        }),
    });
})(GetInternalSquadAccessibleNodesCommand || (exports.GetInternalSquadAccessibleNodesCommand = GetInternalSquadAccessibleNodesCommand = {}));
