"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserAccessibleNodesCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetUserAccessibleNodesCommand;
(function (GetUserAccessibleNodesCommand) {
    GetUserAccessibleNodesCommand.url = api_1.REST_API.USERS.ACCESSIBLE_NODES;
    GetUserAccessibleNodesCommand.TSQ_url = GetUserAccessibleNodesCommand.url(':uuid');
    GetUserAccessibleNodesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.ACCESSIBLE_NODES(':uuid'), 'get', 'Get user accessible nodes');
    GetUserAccessibleNodesCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetUserAccessibleNodesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            userUuid: zod_1.z.string().uuid(),
            activeNodes: zod_1.z.array(zod_1.z.object({
                uuid: zod_1.z.string().uuid(),
                nodeName: zod_1.z.string(),
                countryCode: zod_1.z.string(),
                configProfileUuid: zod_1.z.string().uuid(),
                configProfileName: zod_1.z.string(),
                activeSquads: zod_1.z.array(zod_1.z.object({
                    squadName: zod_1.z.string(),
                    activeInbounds: zod_1.z.array(zod_1.z.string()),
                })),
            })),
        }),
    });
})(GetUserAccessibleNodesCommand || (exports.GetUserAccessibleNodesCommand = GetUserAccessibleNodesCommand = {}));
