"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserAccessibleNodesCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetUserAccessibleNodesCommand;
(function (GetUserAccessibleNodesCommand) {
    GetUserAccessibleNodesCommand.url = api_1.REST_API.USERS.ACCESSIBLE_NODES;
    GetUserAccessibleNodesCommand.TSQ_url = GetUserAccessibleNodesCommand.url(':userId');
    GetUserAccessibleNodesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.ACCESSIBLE_NODES(':userId'), 'get', 'Get user accessible nodes', { scope: 'accessible-nodes', kind: 'read' });
    GetUserAccessibleNodesCommand.RequestParamSchema = zod_1.z.object({
        userId: models_1.numberParamSchema,
    });
    GetUserAccessibleNodesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            userId: zod_1.z.number(),
            activeNodes: zod_1.z.array(zod_1.z.object({
                uuid: zod_1.z.uuid(),
                nodeName: zod_1.z.string(),
                countryCode: zod_1.z.string(),
                configProfileUuid: zod_1.z.uuid(),
                configProfileName: zod_1.z.string(),
                activeSquads: zod_1.z.array(zod_1.z.object({
                    squadName: zod_1.z.string(),
                    activeInbounds: zod_1.z.array(zod_1.z.string()),
                })),
            })),
        }),
    });

})(GetUserAccessibleNodesCommand || (exports.GetUserAccessibleNodesCommand = GetUserAccessibleNodesCommand = {}));
